package usecase

import (
	"context"
	"medioa/internal/storage/models"
)

type IUsecase interface {
	Create(ctx context.Context, userId int64, params *models.SaveRequest) (*models.Response, error)
}
