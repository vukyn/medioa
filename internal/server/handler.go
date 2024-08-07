package server

import (
	"medioa/config"
	initAzBlob "medioa/internal/azblob/init"
	initSecret "medioa/internal/secret/init"
	initStorage "medioa/internal/storage/init"
	"medioa/pkg/log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olahol/melody"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) initHandler(group *gin.RouterGroup) {
	// Init azblob
	azBlob := initAzBlob.NewInit(s.cfg, s.lib)

	// Init secret
	secret := initSecret.NewInit(s.cfg, s.lib)

	// Init storage
	storage := initStorage.NewInit(s.cfg, s.lib, secret, azBlob)
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

func (s *Server) initSwagger() {
	appEnv := s.cfg.App.Environment
	if appEnv == config.APP_ENVIRONMENT_LOCAL {
		s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}

func (s *Server) initStaticFiles() {
	s.router.StaticFS("/index", http.Dir("ui"))
	s.router.GET("/", func(ctx *gin.Context) { ctx.Redirect(http.StatusSeeOther, "/index") })
}

func (s *Server) initSocket() {
	log := log.New("server", "initSocket")
	socket := s.socket
	router := s.router

	router.GET("/ws/:id", func(c *gin.Context) {
		id := c.Param("id")
		socket.HandleRequestWithKeys(c.Writer, c.Request, map[string]any{"id": id})
	})

	router.GET("/ws/start", func(c *gin.Context) {
		id := uuid.New().String()
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	socket.HandleConnect(func(sess *melody.Session) {
		id := sess.Keys["id"].(string)
		log.Info("connected: %s", id)
		s.lib.SocketConn.Set(id, sess)
	})

	socket.HandleError(func(sess *melody.Session, e error) {
		log.Error("error:", e)
	})

	socket.HandleDisconnect(func(sess *melody.Session) {
		id := sess.Keys["id"].(string)
		log.Info("closed: %s", id)
		s.lib.SocketConn.Delete(id)
	})

}

func (s *Server) initCORS() {
	router := s.router
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true

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

	router.Use(cors.New(corsConfig))
}
