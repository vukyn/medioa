package usecase

import (
	"medioa/config"
	azBlobSv "medioa/internal/azblob/service"
	secretSv "medioa/internal/secret/service"
	storageSv "medioa/internal/storage/service"
)

type usecase struct {
	cfg       *config.Config
	storageSv storageSv.IService
	secretSv  secretSv.IService
	azBlobSv  azBlobSv.IService
}

func InitUsecase(cfg *config.Config, storageSv storageSv.IService, secretSv secretSv.IService, azBlobSv azBlobSv.IService) IUsecase {
	return &usecase{
		cfg:       cfg,
		storageSv: storageSv,
		secretSv:  secretSv,
		azBlobSv:  azBlobSv,
	}
}
