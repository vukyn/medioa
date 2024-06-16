package usecase

import (
	"context"
	"fmt"
	"medioa/config"
	"medioa/constants"
	"medioa/internal/storage/models"
	storageSv "medioa/internal/storage/service"
	"medioa/pkg/log"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/zRedShift/mimemagic"
)

type usecase struct {
	cfg       *config.Config
	storageSv storageSv.IService
}

func InitUsecase(cfg *config.Config, storageSv storageSv.IService) IUsecase {
	return &usecase{
		cfg:       cfg,
		storageSv: storageSv,
	}
}

func (u *usecase) Upload(ctx context.Context, userId int64, params *models.UploadRequest) (*models.UploadResponse, error) {
	log := log.New("service", "Upload")

	// sniff mime type
	mimeType, err := sniffMimeType(params.File)
	if err != nil {
		return nil, err
	}

	res, err := u.storageSv.UploadBlob(ctx, params.ToBlobRequest())
	if err != nil {
		log.Error("service.storageSv.UploadBlob", err)
		return nil, err
	}

	// Save to database
	_id := uuid.New().String()
	downloadUrl := fmt.Sprintf("%s%s/%s?token=%s", u.cfg.App.Host, constants.STORAGE_ENDPOINT_DOWNLOAD, _id, res.Token)
	if _, err := u.storageSv.Create(ctx, userId, &models.SaveRequest{
		UUID:        _id,
		Type:        mimeType,
		Token:       res.Token,
		DownloadUrl: downloadUrl,
		Ext:         res.Ext,
	}); err != nil {
		log.Error("service.storageSv.Create", err)
		return nil, err
	}

	return &models.UploadResponse{
		Url:      downloadUrl,
		FileName: _id,
		Token:    res.Token,
		Ext:      res.Ext,
	}, nil
}

func (u *usecase) Download(ctx context.Context, userId int64, params *models.DownloadRequest) (*models.DownloadResponse, error) {
	log := log.New("service", "Download")

	// Validation
	if params.FileName == "" {
		return nil, fmt.Errorf("file name is required")
	}

	if params.Token == "" {
		return nil, fmt.Errorf("token is required")
	}
	// End validation

	// Get file info
	file, err := u.storageSv.GetOne(ctx, &models.RequestParams{
		UUID:  params.FileName,
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

	fileName := file.Token
	if file.Ext != "" {
		fileName += file.Ext
	}
	sas, err := u.storageSv.DownloadSAS(ctx, &models.DownloadSASRequest{
		FileName: fileName,
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
