package server

import (
	"medioa/config"
	initStorage "medioa/internal/storage/init"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

func (s *Server) initCORS(r *gin.Engine) {
	corsConfig := cors.DefaultConfig()

	if len(s.cfg.Cors.AllowOrigins) > 0 {
		corsConfig.AllowOrigins = s.cfg.Cors.AllowOrigins
	} else {
		corsConfig.AllowAllOrigins = true
	}

	if len(s.cfg.Cors.AllowHeaders) > 0 {
		corsConfig.AllowHeaders = s.cfg.Cors.AllowHeaders
	} else {
		corsConfig.AllowHeaders = []string{"*"}
	}
	if len(s.cfg.Cors.AllowMethods) > 0 {
		corsConfig.AllowMethods = s.cfg.Cors.AllowMethods
	} else {
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	}

	r.Use(cors.New(corsConfig))
}
