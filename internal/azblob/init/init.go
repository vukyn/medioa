package init

import (
	"medioa/config"
	"medioa/internal/azblob/service"
	commonModel "medioa/models"
)

type Init struct {
	Service service.IService
}

func NewInit(
	cfg *config.Config,
	lib *commonModel.Lib,
) *Init {
	service := service.InitService(cfg, lib)
	return &Init{
		Service: service,
	}
}
