package usecase

import (
	"context"
	"medioa/internal/storage/models"
	storageSv "medioa/internal/storage/service"
)

type usecase struct {
	storageSv storageSv.IService
}

func InitUsecase(storageSv storageSv.IService) IUsecase {
	return &usecase{
		storageSv: storageSv,
	}
}

func (u *usecase) Create(ctx context.Context, userId int64, params *models.SaveRequest) (*models.Response, error) {
	return u.storageSv.Create(ctx, userId, params)
}
