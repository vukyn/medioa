package server

import (
	"medioa/config"
	initStorage "medioa/internal/storage/init"
	"medioa/pkg/log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	socketio "github.com/googollee/go-socket.io"
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

	s.router.GET("/socket.io/*any", gin.WrapH(s.socket))
	s.router.POST("/socket.io/*any", gin.WrapH(s.socket))

	s.socket.OnConnect("/", func(conn socketio.Conn) error {
		conn.SetContext("")
		log.Info("connected: %s", conn.ID())
		s.lib.SocketConn.Set(conn.ID(), conn)
		return nil
	})

	s.socket.OnError("/", func(s socketio.Conn, e error) {
		log.Error("error:", e)
	})

	s.socket.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		log.Info("closed: %s", reason)
		s.lib.SocketConn.Delete(conn.ID())
	})

}

func (s *Server) initCORS() {
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

	s.router.Use(cors.New(corsConfig))
}
