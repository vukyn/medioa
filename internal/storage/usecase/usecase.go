package usecase

import (
	"context"
	"medioa/config"
	azBlobSv "medioa/internal/azblob/service"
	secretSv "medioa/internal/secret/service"
	storageModel "medioa/internal/storage/models"
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

func (u *usecase) GetFileInfo(ctx context.Context, userId int64, params *storageModel.GetFileInfoRequest) (*storageModel.GetFileInfoResponse, error) {

	// get file info
	file, err := u.getFileById(ctx, params.FileId)
	if err != nil {
		return nil, err
	}

	return &storageModel.GetFileInfoResponse{
		FileId:    file.UUID,
		FileName:  file.FileName,
		FileSize:  file.FileSize,
		HasSecret: file.SecretId != "",
	}, nil
}
