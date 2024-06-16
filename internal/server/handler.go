package server

import (
	"medioa/config"
	initStorage "medioa/internal/storage/init"
	"medioa/pkg/log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
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

func (s *Server) initSwagger() {
	appEnv := s.cfg.App.Environment
	if appEnv == config.APP_ENVIRONMENT_LOCAL {
		s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}

func (s *Server) initStaticFiles() {
	s.router.StaticFS("/index", http.Dir("ui"))
}

func (s *Server) initSocket() {
	log := log.New("server", "initSocket")
	socket := s.socket
	router := s.router

	router.GET("/ws", func(c *gin.Context) {
		socket.HandleRequest(c.Writer, c.Request)
	})

	socket.HandleConnect(func(sess *melody.Session) {
		id := sess.Request.RemoteAddr
		sess.Set("id", id)
		log.Info("connected: %s", id)
		s.lib.SocketConn.Set(id, sess)
		sess.Write([]byte(id)) // Send the socket ID to the client
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
