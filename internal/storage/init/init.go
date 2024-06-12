package init

import (
	"medioa/config"
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

func NewInit(lib *commonModel.Lib, cfg *config.Config) *Init {
	// repository := repository.InitRepo(lib)
	repository := repository.InitMongo(cfg, lib)
	service := service.InitService(cfg, lib, repository)
	usecase := usecase.InitUsecase(cfg, service)
	handler := handler.InitHandler(cfg, lib, usecase)
	return &Init{
		Repository: repository,
		Service:    service,
		Handler:    handler,
		Usecase:    usecase,
	}
}
