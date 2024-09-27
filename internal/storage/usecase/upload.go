package usecase

import (
	"context"
	"fmt"
	azBlobModel "medioa/internal/azblob/models"
	storageModel "medioa/internal/storage/models"
	"medioa/pkg/log"

	"github.com/google/uuid"
)

func (u *usecase) Upload(ctx context.Context, userId int64, params *storageModel.UploadRequest) (*storageModel.UploadResponse, error) {
	log := log.New("usecase", "Upload")

	// validation

	var err error
	var mimeType string
	if params.File != nil {
		// sniff mime type
		mimeType, err = sniffMimeType(params.File)
		if err != nil {
			return nil, err
		}
	}

	// end validation

	var file *azBlobModel.UploadResponse
	if params.URL != "" {
		// upload from url
		file, err = u.azBlobSv.UploadPublicURL(ctx, params.ToURLRequest())
		if err != nil {
			log.Error("usecase.azBlobSv.UploadPublicURL", err)
			return nil, err
		}
	} else if params.File != nil {
		// upload from file
		file, err = u.azBlobSv.UploadPublicBlob(ctx, params.ToBlobRequest())
		if err != nil {
			log.Error("usecase.azBlobSv.UploadPublicBlob", err)
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid upload request")
	}

	var fileName string
	if params.FileName != "" {
		fileName = params.FileName
	} else {
		if params.File != nil {
			fileName = getUploadedFileName1(params.File)
		} else {
			fileName = getUploadedFileName2(params.URL)
		}
	}

	var fileSize int64
	if params.File != nil {
		fileSize = params.File.Size
	}

	// save to database
	fileId := uuid.New().String()
	downloadUrl := getDownloadUrl(u.cfg.App.Host, fileId, file.Token)
	if _, err := u.storageSv.Create(ctx, userId, &storageModel.SaveRequest{
		UUID:        fileId,
		Type:        mimeType,
		Token:       file.Token,
		DownloadUrl: downloadUrl,
		Ext:         file.Ext,
		FileName:    fileName,
		FileSize:    fileSize,
	}); err != nil {
		log.Error("usecase.storageSv.Create", err)
		return nil, err
	}

	return &storageModel.UploadResponse{
		Url:      downloadUrl,
		FileId:   fileId,
		Token:    file.Token,
		Ext:      file.Ext,
		FileName: fileName,
		FileSize: fileSize,
	}, nil
}

func (u *usecase) UploadWithSecret(ctx context.Context, userId int64, params *storageModel.UploadWithSecretRequest) (*storageModel.UploadResponse, error) {
	log := log.New("usecase", "UploadWithSecret")

	// validation

	secret, err := u.verifySecretToken(ctx, params.Secret)
	if err != nil {
		return nil, err
	}

	// sniff mime type
	mimeType, err := sniffMimeType(params.File)
	if err != nil {
		return nil, err
	}

	// end validation

	// upload to private blob
	uploadReq := params.ToBlobRequest(secret.UUID)
	file, err := u.azBlobSv.UploadPrivateBlob(ctx, uploadReq)
	if err != nil {
		log.Error("usecase.azBlobSv.UploadPrivateBlob", err)
		return nil, err
	}

	// Save to database
	fileId := uuid.New().String()
	fileName := map[bool]string{true: params.FileName, false: getUploadedFileName1(params.File)}[params.FileName != ""]
	downloadUrl := getDownloadUrl(u.cfg.App.Host, fileId, file.Token)
	if _, err := u.storageSv.Create(ctx, userId, &storageModel.SaveRequest{
		UUID:        fileId,
		Type:        mimeType,
		Token:       file.Token,
		DownloadUrl: downloadUrl,
		Ext:         file.Ext,
		FileName:    fileName,
		FileSize:    params.File.Size,
		SecretId:    secret.UUID,
	}); err != nil {
		log.Error("usecase.storageSv.Create", err)
		return nil, err
	}

	return &storageModel.UploadResponse{
		Url:      downloadUrl,
		FileId:   fileId,
		Token:    file.Token,
		Ext:      file.Ext,
		FileName: fileName,
		FileSize: params.File.Size,
	}, nil
}

func (u *usecase) UploadChunk(ctx context.Context, userId int64, params *storageModel.UploadChunkRequest) (*storageModel.UploadChunkResponse, error) {
	log := log.New("usecase", "UploadChunk")

	// validation

	// sniff mime type
	mimeType, err := sniffMimeType(params.Chunk)
	if err != nil {
		return nil, err
	}

	// end validation

	// save to database
	var chunkId string
	fileId := params.FileId
	if fileId == "" {
		// upload chunk to azure blob
		file, err := u.azBlobSv.UploadPublicChunk(ctx, params.ToBlobRequest(""))
		if err != nil {
			log.Error("usecase.azBlobSv.UploadPublicChunk", err)
			return nil, err
		}
		chunkId = file.BlockId

		// create new record
		fileId = uuid.New().String()
		fileName := map[bool]string{true: params.FileName, false: getUploadedFileName1(params.Chunk)}[params.FileName != ""]
		downloadUrl := getDownloadUrl(u.cfg.App.Host, fileId, file.Token)
		if _, err := u.storageSv.Create(ctx, userId, &storageModel.SaveRequest{
			UUID:        fileId,
			Type:        mimeType,
			Token:       file.Token,
			DownloadUrl: downloadUrl,
			Ext:         file.Ext,
			FileName:    fileName,
			ChunkIds:    &[]string{file.BlockId},
		}); err != nil {
			log.Error("usecase.storageSv.Create", err)
			return nil, err
		}
	} else {
		// get file info
		file, err := u.getFileById(ctx, params.FileId)
		if err != nil {
			return nil, err
		}

		// upload chunk to azure blob
		_file, err := u.azBlobSv.UploadPublicChunk(ctx, params.ToBlobRequest(file.Token))
		if err != nil {
			log.Error("usecase.azBlobSv.UploadPublicChunk", err)
			return nil, err
		}
		chunkId = _file.BlockId

		// append chunk ids
		chunkIds := append(file.ChunkIds, _file.BlockId)
		if _, err := u.storageSv.Update(ctx, userId, &storageModel.SaveRequest{
			UUID:     params.FileId,
			ChunkIds: &chunkIds,
		}); err != nil {
			log.Error("usecase.storageSv.Update", err)
			return nil, err
		}
	}

	return &storageModel.UploadChunkResponse{
		ChunkId: chunkId,
		FileId:  fileId,
	}, nil
}

func (u *usecase) CommitChunk(ctx context.Context, userId int64, params *storageModel.CommitChunkRequest) (*storageModel.CommitChunkResponse, error) {
	log := log.New("usecase", "CommitChunk")

	// validation

	// get file info
	file, err := u.getFileById(ctx, params.FileId)
	if err != nil {
		return nil, err
	}

	// end validation

	res, err := u.azBlobSv.CommitPublicChunk(ctx, &azBlobModel.CommitChunkRequest{
		SessionId: params.SessionId,
		Token:     file.Token,
		FileName:  file.FileName,
		BlockIds:  file.ChunkIds,
	})
	if err != nil {
		log.Error("usecase.azBlobSv.CommitPublicChunk", err)
		return nil, err
	}

	// update file info
	if _, err := u.storageSv.Update(ctx, userId, &storageModel.SaveRequest{
		UUID:        params.FileId,
		FileSize:    res.FileSize,
		TotalChunks: res.TotalBlock,
		ChunkIds:    &[]string{},
	}); err != nil {
		log.Error("usecase.storageSv.Update", err)
		return nil, err
	}

	return &storageModel.CommitChunkResponse{
		Url:      file.DownloadUrl,
		FileId:   file.UUID,
		Token:    file.Token,
		Ext:      file.Ext,
		FileName: file.FileName,
		FileSize: res.FileSize,
	}, nil
}

func (u *usecase) UploadChunkWithSecret(ctx context.Context, userId int64, params *storageModel.UploadChunkWithSecretRequest) (*storageModel.UploadChunkResponse, error) {
	log := log.New("usecase", "UploadChunkWithSecret")

	// validation

	// get secret info
	secret, err := u.verifySecretToken(ctx, params.Secret)
	if err != nil {
		return nil, err
	}

	// sniff mime type
	mimeType, err := sniffMimeType(params.Chunk)
	if err != nil {
		return nil, err
	}

	// end validation

	// save to database
	var chunkId string
	fileId := params.FileId
	if fileId == "" {
		// upload chunk to azure blob
		file, err := u.azBlobSv.UploadPrivateChunk(ctx, params.ToBlobRequest(secret.UUID, ""))
		if err != nil {
			log.Error("usecase.azBlobSv.UploadPrivateChunk", err)
			return nil, err
		}
		chunkId = file.BlockId

		// Create new record
		fileId = uuid.New().String()
		fileName := map[bool]string{true: params.FileName, false: getUploadedFileName1(params.Chunk)}[params.FileName != ""]
		downloadUrl := getDownloadUrl(u.cfg.App.Host, fileId, file.Token)
		if _, err := u.storageSv.Create(ctx, userId, &storageModel.SaveRequest{
			UUID:        fileId,
			Type:        mimeType,
			Token:       file.Token,
			DownloadUrl: downloadUrl,
			Ext:         file.Ext,
			FileName:    fileName,
			ChunkIds:    &[]string{file.BlockId},
			SecretId:    secret.UUID,
		}); err != nil {
			log.Error("usecase.storageSv.Create", err)
			return nil, err
		}
	} else {
		// get file info
		file, err := u.getFileById(ctx, params.FileId)
		if err != nil {
			return nil, err
		}

		// upload chunk to azure blob
		_file, err := u.azBlobSv.UploadPrivateChunk(ctx, params.ToBlobRequest(secret.UUID, file.Token))
		if err != nil {
			log.Error("usecase.azBlobSv.UploadPrivateChunk", err)
			return nil, err
		}
		chunkId = _file.BlockId

		// append chunk ids
		chunkIds := append(file.ChunkIds, _file.BlockId)
		if _, err := u.storageSv.Update(ctx, userId, &storageModel.SaveRequest{
			UUID:     params.FileId,
			ChunkIds: &chunkIds,
		}); err != nil {
			log.Error("usecase.storageSv.Update", err)
			return nil, err
		}
	}

	return &storageModel.UploadChunkResponse{
		ChunkId: chunkId,
		FileId:  fileId,
	}, nil
}

func (u *usecase) CommitChunkWithSecret(ctx context.Context, userId int64, params *storageModel.CommitChunkRequest) (*storageModel.CommitChunkResponse, error) {
	log := log.New("usecase", "CommitChunkWithSecret")

	// validation

	// get file info
	file, err := u.getFileById(ctx, params.FileId)
	if err != nil {
		return nil, err
	}

	// get secret info
	secret, err := u.verifySecretToken(ctx, params.Secret)
	if err != nil {
		return nil, err
	}

	// end validation

	res, err := u.azBlobSv.CommitPrivateChunk(ctx, &azBlobModel.CommitChunkRequest{
		SessionId: params.SessionId,
		SecretId:  secret.UUID,
		Token:     file.Token,
		FileName:  file.FileName,
		BlockIds:  file.ChunkIds,
	})
	if err != nil {
		log.Error("usecase.azBlobSv.CommitPrivateChunk", err)
		return nil, err
	}

	// update file info
	if _, err := u.storageSv.Update(ctx, userId, &storageModel.SaveRequest{
		UUID:        params.FileId,
		FileSize:    res.FileSize,
		TotalChunks: res.TotalBlock,
		ChunkIds:    &[]string{},
	}); err != nil {
		log.Error("usecase.storageSv.Update", err)
		return nil, err
	}

	return &storageModel.CommitChunkResponse{
		Url:      file.DownloadUrl,
		FileId:   file.UUID,
		Token:    file.Token,
		Ext:      file.Ext,
		FileName: file.FileName,
		FileSize: res.FileSize,
	}, nil
}
