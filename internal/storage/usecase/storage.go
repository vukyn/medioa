package usecase

import (
	"context"
	"fmt"
	"medioa/constants"
	azBlobModel "medioa/internal/azblob/models"
	secretModel "medioa/internal/secret/models"
	storageModel "medioa/internal/storage/models"
	"medioa/pkg/log"
	"medioa/pkg/xtype"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/zRedShift/mimemagic"
)

func (u *usecase) Upload(ctx context.Context, userId int64, params *storageModel.UploadFileRequest) (*storageModel.UploadResponse, error) {
	log := log.New("usecase", "Upload")

	// sniff mime type
	mimeType, err := sniffMimeType(params.File)
	if err != nil {
		return nil, err
	}

	file, err := u.azBlobSv.UploadPublicBlob(ctx, params.ToBlobRequest())
	if err != nil {
		log.Error("usecase.storageSv.UploadPublicBlob", err)
		return nil, err
	}

	fileName := params.FileName
	if fileName == "" {
		fileName = getUploadedFileName(params.File)
	}

	// Save to database
	fileId := uuid.New().String()
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

func (u *usecase) UploadChunk(ctx context.Context, userId int64, params *storageModel.UploadChunkRequest) (*storageModel.UploadChunkResponse, error) {
	log := log.New("usecase", "UploadChunk")

	// sniff mime type
	mimeType, err := sniffMimeType(params.Chunk)
	if err != nil {
		return nil, err
	}

	fileName := params.FileName
	if fileName == "" {
		fileName = getUploadedFileName(params.Chunk)
	}

	// Save to database
	var chunkId string
	fileId := params.FileId
	if fileId == "" {
		// upload chunk to azure blob
		file, err := u.azBlobSv.UploadPublicChunk(ctx, params.ToBlobRequest(""))
		if err != nil {
			log.Error("usecase.storageSv.UploadPublicChunk", err)
			return nil, err
		}
		chunkId = file.BlockId

		// Create new record
		fileId = uuid.New().String()
		filePath := strings.ReplaceAll(constants.STORAGE_ENDPOINT_DOWNLOAD, ":file_id", fileId)
		downloadUrl := fmt.Sprintf("%s/api/v1%s?token=%s", u.cfg.App.Host, filePath, file.Token)
		if _, err := u.storageSv.Create(ctx, userId, &storageModel.SaveRequest{
			UUID:        fileId,
			Type:        mimeType,
			Token:       file.Token,
			DownloadUrl: downloadUrl,
			Ext:         file.Ext,
			FileName:    fileName,
			FileSize:    params.TotalChunks,
			ChunkIds:    []string{file.BlockId},
		}); err != nil {
			log.Error("usecase.storageSv.Create", err)
			return nil, err
		}
	} else {
		// Get file info
		storage, err := u.storageSv.GetOne(ctx, &storageModel.RequestParams{
			UUID: params.FileId,
		})
		if err != nil {
			log.Error("usecase.storageSv.GetOne", err)
			return nil, err
		}
		if storage == nil {
			return nil, fmt.Errorf("file not found")
		}

		// upload chunk to azure blob
		file, err := u.azBlobSv.UploadPublicChunk(ctx, params.ToBlobRequest(storage.Token))
		if err != nil {
			log.Error("usecase.storageSv.UploadPublicChunk", err)
			return nil, err
		}
		chunkId = file.BlockId

		// Append chunk ids
		if _, err := u.storageSv.Update(ctx, userId, &storageModel.SaveRequest{
			UUID:     params.FileId,
			ChunkIds: append(storage.ChunkIds, file.BlockId),
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

	// Validation
	if params.FileId == "" {
		return nil, fmt.Errorf("file id is required")
	}

	// Get file info
	storage, err := u.storageSv.GetOne(ctx, &storageModel.RequestParams{
		UUID: params.FileId,
	})
	if err != nil {
		log.Error("usecase.storageSv.GetOne", err)
		return nil, err
	}
	if storage == nil {
		return nil, fmt.Errorf("file not found")
	}

	if err := u.azBlobSv.CommitPublicChunk(ctx, &azBlobModel.CommitChunkRequest{
		SessionId: params.SessionId,
		Token:     storage.Token,
		FileName:  storage.FileName,
		BlockIds:  storage.ChunkIds,
	}); err != nil {
		log.Error("usecase.storageSv.CommitPublicChunk", err)
		return nil, err
	}

	return &storageModel.CommitChunkResponse{
		Url:      storage.DownloadUrl,
		FileId:   storage.UUID,
		Token:    storage.Token,
		Ext:      storage.Ext,
		FileName: storage.FileName,
		FileSize: storage.FileSize,
	}, nil
}

func (u *usecase) UploadWithSecret(ctx context.Context, userId int64, params *storageModel.UploadWithSecretRequest) (*storageModel.UploadResponse, error) {
	log := log.New("usecase", "UploadWithSecret")

	secret, err := u.secretSv.GetOne(ctx, &secretModel.RequestParams{
		AccessToken: params.Secret,
	})
	if err != nil {
		log.Error("usecase.secretSv.GetOne", err)
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("secret token is invalid")
	}

	// sniff mime type
	mimeType, err := sniffMimeType(params.File)
	if err != nil {
		return nil, err
	}

	uploadReq := params.ToBlobRequest()
	uploadReq.SecretId = secret.UUID
	file, err := u.azBlobSv.UploadPrivateBlob(ctx, uploadReq)
	if err != nil {
		log.Error("usecase.storageSv.UploadPrivateBlob", err)
		return nil, err
	}

	fileName := params.FileName
	if fileName == "" {
		fileName = getUploadedFileName(params.File)
	}

	// Save to database
	fileId := uuid.New().String()
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

func (u *usecase) Download(ctx context.Context, userId int64, params *storageModel.DownloadRequest) (*storageModel.DownloadResponse, error) {
	log := log.New("usecase", "Download")

	// Validation
	if params.FileId == "" {
		return nil, fmt.Errorf("file id is required")
	}

	if params.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	// Get file info
	file, err := u.storageSv.GetOne(ctx, &storageModel.RequestParams{
		UUID:  params.FileId,
		Token: params.Token,
	})
	if err != nil {
		log.Error("usecase.storageSv.GetOne", err)
		return nil, err
	}
	if file == nil {
		return nil, fmt.Errorf("file not found")
	}

	// Check permission
	if file.SecretId != "" {
		return nil, fmt.Errorf("permission denied")
	}
	if file.CreatedBy != userId {
		return nil, fmt.Errorf("permission denied")
	}
	// End validation

	downloadFileName := file.Token
	if file.Ext != "" {
		downloadFileName += file.Ext
	}
	downloadFileName = path.Join("public", downloadFileName)
	sas, err := u.azBlobSv.DownloadSAS(ctx, &azBlobModel.DownloadSASRequest{
		FileName: downloadFileName,
	})
	if err != nil {
		log.Error("usecase.storageSv.DownloadSAS", err)
		return nil, err
	}

	return &storageModel.DownloadResponse{
		Url: sas.Url,
	}, nil
}

func (u *usecase) DownloadWithSecret(ctx context.Context, userId int64, params *storageModel.DownloadWithSecretRequest) (*storageModel.DownloadResponse, error) {
	log := log.New("usecase", "DownloadWithSecret")

	// Validation
	if params.FileId == "" {
		return nil, fmt.Errorf("file id is required")
	}

	if params.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	// Get file info
	file, err := u.storageSv.GetOne(ctx, &storageModel.RequestParams{
		UUID:  params.FileId,
		Token: params.Token,
	})
	if err != nil {
		log.Error("usecase.storageSv.GetOne", err)
		return nil, err
	}
	if file == nil {
		return nil, fmt.Errorf("file not found")
	}

	// Get secret info
	secret, err := u.secretSv.GetOne(ctx, &secretModel.RequestParams{
		AccessToken: params.Secret,
	})
	if err != nil {
		log.Error("usecase.secretSv.GetOne", err)
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("secret token is invalid")
	}

	// Check permission
	if file.SecretId != secret.UUID {
		return nil, fmt.Errorf("permission denied")
	}
	if file.CreatedBy != userId {
		return nil, fmt.Errorf("permission denied")
	}
	// End validation

	downloadFileName := file.Token
	if file.Ext != "" {
		downloadFileName += file.Ext
	}
	downloadFileName = path.Join("private", file.SecretId, downloadFileName)
	sas, err := u.azBlobSv.DownloadSAS(ctx, &azBlobModel.DownloadSASRequest{
		FileName: downloadFileName,
	})
	if err != nil {
		log.Error("usecase.storageSv.DownloadSAS", err)
		return nil, err
	}

	return &storageModel.DownloadResponse{
		Url: sas.Url,
	}, nil
}

func sniffMimeType(file *xtype.File) (string, error) {
	log := log.New("usecase", "sniffMimeType")

	reader, err := file.Open()
	if err != nil {
		log.Error("file.Open", err)
		return "", err
	}
	defer reader.Close()

	mimeType, err := mimemagic.MatchReader(reader, file.Filename)
	if err != nil {
		log.Error("mimemagic.MatchReader", err)
		return "", err
	}

	return mimeType.MediaType(), nil
}

func getUploadedFileName(file *xtype.File) string {
	var fileName string
	ext := path.Ext(file.Filename)
	fileName = strings.ReplaceAll(file.Filename, ext, "")
	if fileName == "" {
		fileName = file.Filename
	}
	return fileName
}
