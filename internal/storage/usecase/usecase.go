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

	// sniff mime type
	mimeType, err := sniffMimeType(params.File)
	if err != nil {
		return nil, err
	}

	res, err := u.storageSv.Upload(ctx, params)
	if err != nil {
		return nil, err
	}

	// Save to database
	_id := uuid.New()
	downloadUrl := fmt.Sprintf("%s%s/%s?token=%s", u.cfg.App.Host, constants.STORAGE_ENDPOINT_DOWNLOAD, _id.String(), res.Token)
	if _, err := u.storageSv.Create(ctx, userId, &models.SaveRequest{
		UUID:        _id,
		Type:        mimeType,
		Token:       res.Token,
		DownloadUrl: downloadUrl,
	}); err != nil {
		return nil, err
	}

	return &models.UploadResponse{
		Url:      downloadUrl,
		FileName: _id.String(),
		Token:    res.Token,
		Ext:      res.Ext,
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
