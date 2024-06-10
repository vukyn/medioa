package server

import (
	"context"
	"io"
	"medioa/config"
	"medioa/models"
	"medioa/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	lib *models.Lib
	cfg *config.Config
}

func New(ctx context.Context, cfg *config.Config) *Server {
	mongoClient, err := initMongo(ctx, cfg)
	if err != nil {
		panic(err)
	}

	lib := &models.Lib{
		Mongo: mongoClient,
	}
	return &Server{
		cfg: cfg,
		lib: lib,
	}
}

func (s *Server) Start(ctx context.Context) {
	log := log.New("server", "Start")

	port := s.cfg.AppConfig.Port
	appEnv := s.cfg.AppConfig.Environment

	if appEnv == config.APP_ENVIRONMENT_PROD {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}

	r := initGin()
	v1 := r.Group("/api/v1")
	s.initHealthCheck(v1)
	s.initHandler(v1)

	if appEnv == config.APP_ENVIRONMENT_LOCAL {
		// s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	log.Info("started api successfully with port: %v", port)
	http.ListenAndServe(port, r)
}

func (s *Server) Stop(ctx context.Context) {
	log := log.New("server", "Stop")

	log.Info("stopping api")
	if err := s.lib.Mongo.Disconnect(ctx); err != nil {
		log.Error("failed to disconnect mongo", err)
	}
}

func initGin() *gin.Engine {
	log := log.New("server", "initGin")

	r := gin.Default()
	r.Use(gin.Recovery())
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Debug("endpoint %v %v %v", httpMethod, absolutePath, handlerName)
	}
	return r
}

func (s *Server) initHealthCheck(group *gin.RouterGroup) {
	pingHandler := func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]any{
			"version":   s.cfg.AppConfig.Version,
			"client_ip": ctx.ClientIP(),
		})
	}

	group.GET("/health-check", pingHandler)
}

func initMongo(ctx context.Context, cfg *config.Config) (*mongo.Client, error) {
	log := log.New("server", "initMongo")

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.MongoConfig.URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Error("mongo.Connect", err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error("client.Ping", err)
		return nil, err
	}
	log.Info("connected to mongo successfully")

	return client, nil
}
