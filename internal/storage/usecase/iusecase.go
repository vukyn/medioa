package usecase

import (
	"context"
	"medioa/internal/storage/models"
)

type IUsecase interface {
	Upload(ctx context.Context, userId int64, params *models.UploadRequest) (*models.UploadResponse, error)
	UploadWithSecret(ctx context.Context, userId int64, params *models.UploadWithSecretRequest) (*models.UploadResponse, error)
	Download(ctx context.Context, userId int64, params *models.DownloadRequest) (*models.DownloadResponse, error)
	DownloadWithSecret(ctx context.Context, userId int64, params *models.DownloadWithSecretRequest) (*models.DownloadResponse, error)
	CreateSecret(ctx context.Context, userId int64, params *models.CreateSecretRequest) (*models.CreateSecretResponse, error)
	RetrieveSecret(ctx context.Context, userId int64, params *models.RetrieveSecretRequest) (*models.RetrieveSecretResponse, error)
	ResetPinCode(ctx context.Context, userId int64, params *models.ResetPinCodeRequest) (int64, error)
}
