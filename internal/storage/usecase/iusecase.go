package usecase

import (
	"context"
	"medioa/internal/storage/models"
)

type IUsecase interface {
	Upload(ctx context.Context, userId int64, params *models.UploadFileRequest) (*models.UploadResponse, error)
	UploadChunk(ctx context.Context, userId int64, params *models.UploadChunkRequest) (*models.UploadChunkResponse, error)
	CommitChunk(ctx context.Context, userId int64, params *models.CommitChunkRequest) (*models.CommitChunkResponse, error)
	UploadWithSecret(ctx context.Context, userId int64, params *models.UploadWithSecretRequest) (*models.UploadResponse, error)
	UploadChunkWithSecret(ctx context.Context, userId int64, params *models.UploadChunkWithSecretRequest) (*models.UploadChunkResponse, error)
	CommitChunkWithSecret(ctx context.Context, userId int64, params *models.CommitChunkRequest) (*models.CommitChunkResponse, error)
	Download(ctx context.Context, userId int64, params *models.DownloadRequest) (*models.DownloadResponse, error)
	DownloadWithSecret(ctx context.Context, userId int64, params *models.DownloadWithSecretRequest) (*models.DownloadResponse, error)
	CreateSecret(ctx context.Context, userId int64, params *models.CreateSecretRequest) (*models.CreateSecretResponse, error)
	RetrieveSecret(ctx context.Context, userId int64, params *models.RetrieveSecretRequest) (*models.RetrieveSecretResponse, error)
	ResetPinCode(ctx context.Context, userId int64, params *models.ResetPinCodeRequest) (int64, error)
}
