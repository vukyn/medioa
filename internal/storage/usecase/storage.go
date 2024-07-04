package usecase

import (
	"context"
	"fmt"
	"medioa/constants"
	secretModel "medioa/internal/secret/models"
	storageModel "medioa/internal/storage/models"
	"medioa/pkg/log"
	"mime/multipart"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/zRedShift/mimemagic"
)

func (u *usecase) Upload(ctx context.Context, userId int64, params *storageModel.UploadRequest) (*storageModel.UploadResponse, error) {
	log := log.New("service", "Upload")

	// sniff mime type
	mimeType, err := sniffMimeType(params.File)
	if err != nil {
		return nil, err
	}

	file, err := u.storageSv.UploadPublicBlob(ctx, params.ToBlobRequest())
	if err != nil {
		log.Error("service.storageSv.UploadPublicBlob", err)
		return nil, err
	}

	fileName := params.FileName
	if fileName == "" {
		ext := path.Ext(params.File.Filename)
		fileName = strings.ReplaceAll(params.File.Filename, ext, "")
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
		log.Error("service.storageSv.Create", err)
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
	log := log.New("service", "UploadWithSecret")

	secret, err := u.secretSv.GetOne(ctx, &secretModel.RequestParams{
		AccessToken: params.Secret,
	})
	if err != nil {
		log.Error("service.secretSv.GetOne", err)
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
	file, err := u.storageSv.UploadPrivateBlob(ctx, uploadReq)
	if err != nil {
		log.Error("service.storageSv.UploadPrivateBlob", err)
		return nil, err
	}

	fileName := params.FileName
	if fileName == "" {
		ext := path.Ext(params.File.Filename)
		fileName = strings.ReplaceAll(params.File.Filename, ext, "")
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
		log.Error("service.storageSv.Create", err)
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
	log := log.New("service", "Download")

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
		log.Error("service.storageSv.GetOne", err)
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
	sas, err := u.storageSv.DownloadSAS(ctx, &storageModel.DownloadSASRequest{
		FileName: downloadFileName,
	})
	if err != nil {
		log.Error("service.storageSv.DownloadSAS", err)
		return nil, err
	}

	return &storageModel.DownloadResponse{
		Url: sas.Url,
	}, nil
}

func (u *usecase) DownloadWithSecret(ctx context.Context, userId int64, params *storageModel.DownloadWithSecretRequest) (*storageModel.DownloadResponse, error) {
	log := log.New("service", "DownloadWithSecret")

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
		log.Error("service.storageSv.GetOne", err)
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
		log.Error("service.secretSv.GetOne", err)
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
	sas, err := u.storageSv.DownloadSAS(ctx, &storageModel.DownloadSASRequest{
		FileName: downloadFileName,
	})
	if err != nil {
		log.Error("service.storageSv.DownloadSAS", err)
		return nil, err
	}

	return &storageModel.DownloadResponse{
		Url: sas.Url,
	}, nil
}

func sniffMimeType(file *multipart.FileHeader) (string, error) {
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
