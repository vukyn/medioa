package main

import (
	"context"
	"medioa/config"
	"medioa/internal/server"
	"time"

	"medioa/pkg/graceful"
	"medioa/pkg/log"
)

//	@title			Medioa API
//	@version		1.0
//	@description	Medioa REST API (with gin-gonic).

//	@contact.name	Vũ Kỳ
//	@contact.url	github.com/vukyn
//	@contact.email	vukynpro@gmail.com

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

// @BasePath	/api/v1
func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log.Init(cfg.Log)
	server := server.New(ctx, cfg)
	go func() {
		server.Start(ctx)
	}()
	defer server.Stop(ctx)

	graceful.ShutDownSlowly(time.Duration(cfg.App.ShutdownTimeout) * time.Second)
}
