package usecase

import (
	"medioa/config"
	secretSv "medioa/internal/secret/service"
	storageSv "medioa/internal/storage/service"
)

type usecase struct {
	cfg       *config.Config
	storageSv storageSv.IService
	secretSv  secretSv.IService
}

func InitUsecase(cfg *config.Config, storageSv storageSv.IService, secretSv secretSv.IService) IUsecase {
	return &usecase{
		cfg:       cfg,
		storageSv: storageSv,
		secretSv:  secretSv,
	}
}
