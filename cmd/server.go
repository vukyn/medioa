package main

import (
	"context"
	"medioa/config"
	"medioa/internal/server"
	"time"

	"medioa/pkg/graceful"
	"medioa/pkg/log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log.Init(cfg.LogConfig)
	server := server.New(ctx, cfg)
	go func() {
		server.Start(ctx)
	}()
	defer server.Stop(ctx)

	graceful.ShutDownSlowly(time.Duration(cfg.AppConfig.ShutdownTimeout) * time.Second)
}
