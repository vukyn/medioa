package usecase

import (
	"context"
	"fmt"
	"medioa/constants"
	azBlobModel "medioa/internal/azblob/models"
	storageModel "medioa/internal/storage/models"
	"medioa/pkg/log"
	"strings"

	"github.com/google/uuid"
)

func (u *usecase) Upload(ctx context.Context, userId int64, params *storageModel.UploadFileRequest) (*storageModel.UploadResponse, error) {
	log := log.New("usecase", "Upload")

	// validation

	// sniff mime type
	mimeType, err := sniffMimeType(params.File)
	if err != nil {
		return nil, err
	}

	// end validation

	file, err := u.azBlobSv.UploadPublicBlob(ctx, params.ToBlobRequest())
	if err != nil {
		log.Error("usecase.azBlobSv.UploadPublicBlob", err)
		return nil, err
	}

	// save to database
	fileId := uuid.New().String()
	fileName := map[bool]string{true: params.FileName, false: getUploadedFileName(params.File)}[params.FileName != ""]
	filePath := strings.ReplaceAll(constants.STORAGE_ENDPOINT_DOWNLOAD, ":file_id", fileId)
	downloadUrl := fmt.Sprintf("%s/api/v1%s?token=%s", u.cfg.App.Host, filePath, file.Token)
	if _, err := u.storageSv.Create(ctx, userId, &storageModel.SaveRequest{
		UUID:        fileId,
		Type:        mimeType,
		Token:       file.Token,
		DownloadUrl: downloadUrl,
		Ext:         file.Ext,
		FileName:    fileName,
		FileSize:    params.File.Size,
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
	fileName := map[bool]string{true: params.FileName, false: getUploadedFileName(params.File)}[params.FileName != ""]
	filePath := strings.ReplaceAll(constants.STORAGE_ENDPOINT_DOWNLOAD_WITH_SECRET, ":file_id", fileId)
	downloadUrl := fmt.Sprintf("%s/api/v1%s?token=%s", u.cfg.App.Host, filePath, file.Token)
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
		fileName := map[bool]string{true: params.FileName, false: getUploadedFileName(params.Chunk)}[params.FileName != ""]
		filePath := strings.ReplaceAll(constants.STORAGE_ENDPOINT_DOWNLOAD, ":file_id", fileId)
		downloadUrl := fmt.Sprintf("%s/api/v1%s?token=%s", u.cfg.App.Host, filePath, file.Token)
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
		fileName := map[bool]string{true: params.FileName, false: getUploadedFileName(params.Chunk)}[params.FileName != ""]
		filePath := strings.ReplaceAll(constants.STORAGE_ENDPOINT_DOWNLOAD, ":file_id", fileId)
		downloadUrl := fmt.Sprintf("%s/api/v1%s?token=%s", u.cfg.App.Host, filePath, file.Token)
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
