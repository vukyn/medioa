package usecase

import (
	"context"
	"fmt"
	"medioa/constants"
	"medioa/internal/storage/models"
	"medioa/pkg/log"
	"mime/multipart"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/zRedShift/mimemagic"
)

func (u *usecase) Upload(ctx context.Context, userId int64, params *models.UploadRequest) (*models.UploadResponse, error) {
	log := log.New("service", "Upload")

	// sniff mime type
	mimeType, err := sniffMimeType(params.File)
	if err != nil {
		return nil, err
	}

	file, err := u.storageSv.UploadBlob(ctx, params.ToBlobRequest())
	if err != nil {
		log.Error("service.storageSv.UploadBlob", err)
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
	if _, err := u.storageSv.Create(ctx, userId, &models.SaveRequest{
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

	return &models.UploadResponse{
		Url:      downloadUrl,
		FileId:   fileId,
		Token:    file.Token,
		Ext:      file.Ext,
		FileName: fileName,
		FileSize: params.File.Size,
	}, nil
}

func (u *usecase) UploadWithSecret(ctx context.Context, userId int64, params *models.UploadWithSecretRequest) (*models.UploadResponse, error) {
	log := log.New("service", "UploadWithSecret")

	// sniff mime type
	mimeType, err := sniffMimeType(params.File)
	if err != nil {
		return nil, err
	}

	file, err := u.storageSv.UploadBlob(ctx, params.ToBlobRequest())
	if err != nil {
		log.Error("service.storageSv.UploadBlob", err)
		return nil, err
	}

	fileName := params.FileName
	if fileName == "" {
		fileName = params.File.Filename
	}

	// Save to database
	fileId := uuid.New().String()
	filePath := strings.ReplaceAll(constants.STORAGE_ENDPOINT_DOWNLOAD, ":file_id", fileId)
	downloadUrl := fmt.Sprintf("%s/api/v1%s?token=%s", u.cfg.App.Host, filePath, file.Token)
	if _, err := u.storageSv.Create(ctx, userId, &models.SaveRequest{
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

	return &models.UploadResponse{
		Url:      downloadUrl,
		FileId:   fileId,
		Token:    file.Token,
		Ext:      file.Ext,
		FileName: fileName,
		FileSize: params.File.Size,
	}, nil
}

func (u *usecase) Download(ctx context.Context, userId int64, params *models.DownloadRequest) (*models.DownloadResponse, error) {
	log := log.New("service", "Download")

	// Validation
	if params.FileId == "" {
		return nil, fmt.Errorf("file id is required")
	}

	if params.Token == "" {
		return nil, fmt.Errorf("token is required")
	}
	// End validation

	// Get file info
	file, err := u.storageSv.GetOne(ctx, &models.RequestParams{
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
	if file.CreatedBy != userId {
		return nil, fmt.Errorf("permission denied")
	}

	downloadFileName := file.Token
	if file.Ext != "" {
		downloadFileName += file.Ext
	}
	sas, err := u.storageSv.DownloadSAS(ctx, &models.DownloadSASRequest{
		FileName: downloadFileName,
	})
	if err != nil {
		log.Error("service.storageSv.DownloadSAS", err)
		return nil, err
	}

	return &models.DownloadResponse{
		Url: sas.Url,
	}, nil
}

func (u *usecase) DownloadWithSecret(ctx context.Context, userId int64, params *models.DownloadWithSecretRequest) (*models.DownloadResponse, error) {
	log := log.New("service", "DownloadWithSecret")

	// Validation
	if params.FileId == "" {
		return nil, fmt.Errorf("file id is required")
	}

	if params.Token == "" {
		return nil, fmt.Errorf("token is required")
	}
	// End validation

	// Get file info
	file, err := u.storageSv.GetOne(ctx, &models.RequestParams{
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
	if file.CreatedBy != userId {
		return nil, fmt.Errorf("permission denied")
	}

	downloadFileName := file.Token
	if file.Ext != "" {
		downloadFileName += file.Ext
	}
	sas, err := u.storageSv.DownloadSAS(ctx, &models.DownloadSASRequest{
		FileName: downloadFileName,
	})
	if err != nil {
		log.Error("service.storageSv.DownloadSAS", err)
		return nil, err
	}

	return &models.DownloadResponse{
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
