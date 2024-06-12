package server

import (
	initStorage "medioa/internal/storage/init"
	"net/http"
	"medioa/config"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerfiles "github.com/swaggo/files"
)

func (s *Server) initHandler(group *gin.RouterGroup) {
	// Init storage
	storage := initStorage.NewInit(s.lib, s.cfg)
	storage.Handler.MapRoutes(group)

}

func (s *Server) initHealthCheck(group *gin.RouterGroup) {
	pingHandler := func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]any{
			"version":   s.cfg.App.Version,
			"client_ip": ctx.ClientIP(),
		})
	}

	group.GET("/health-check", pingHandler)
}

func (s *Server) initSwagger(r *gin.Engine) {
	appEnv := s.cfg.App.Environment
	if appEnv == config.APP_ENVIRONMENT_LOCAL {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
