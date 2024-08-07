package init

import (
	"medioa/config"
	initAzBlob "medioa/internal/azblob/init"
	initSecret "medioa/internal/secret/init"
	"medioa/internal/storage/handler"
	"medioa/internal/storage/repository"
	"medioa/internal/storage/service"
	"medioa/internal/storage/usecase"
	commonModel "medioa/models"
)

type Init struct {
	Repository repository.IRepository
	Service    service.IService
	Handler    handler.IHandler
	Usecase    usecase.IUsecase
}

func NewInit(
	cfg *config.Config,
	lib *commonModel.Lib,
	initSecret *initSecret.Init,
	initAzBlob *initAzBlob.Init,
) *Init {
	// repository := repository.InitRepo(lib)
	repository := repository.InitMongo(cfg, lib)
	service := service.InitService(cfg, lib, repository)
	usecase := usecase.InitUsecase(cfg, service, initSecret.Service, initAzBlob.Service)
	handler := handler.InitHandler(cfg, lib, usecase)
	return &Init{
		Repository: repository,
		Service:    service,
		Handler:    handler,
		Usecase:    usecase,
	}
}
