package usecase

import (
	"context"
	"medioa/internal/storage/models"
)

type IUsecase interface {
	Upload(ctx context.Context, userId int64, params *models.UploadRequest) (*models.UploadResponse, error)
}
