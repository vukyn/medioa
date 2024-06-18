package init

import (
	"medioa/config"
	"medioa/internal/secret/repository"
	"medioa/internal/secret/service"
	commonModel "medioa/models"
)

type Init struct {
	Repository repository.IRepository
	Service    service.IService
}

func NewInit(
	cfg *config.Config,
	lib *commonModel.Lib,
) *Init {
	// repository := repository.InitRepo(lib)
	repository := repository.InitMongo(cfg, lib)
	service := service.InitService(cfg, lib, repository)
	return &Init{
		Repository: repository,
		Service:    service,
	}
}
