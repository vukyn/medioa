package init

import (
	"medioa/config"
	"medioa/internal/share/handler"
	initStorage "medioa/internal/storage/init"
	commonModel "medioa/models"
)

type Init struct {
	Handler handler.IHandler
}

func NewInit(
	cfg *config.Config,
	lib *commonModel.Lib,
	initSecret *initStorage.Init,
) *Init {
	handler := handler.InitHandler(cfg, lib, initSecret)
	return &Init{
		Handler: handler,
	}
}
