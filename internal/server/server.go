package server

import (
	"context"
	"fmt"
	"io"
	"medioa/config"
	"medioa/models"
	"medioa/pkg/log"
	"medioa/pkg/network"

	_ "medioa/docs"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	lib    *models.Lib
	cfg    *config.Config
	router *gin.Engine
	socket *melody.Melody
}

func New(ctx context.Context, cfg *config.Config) *Server {
	mongoCli, err := initMongo(ctx, cfg)
	if err != nil {
		panic(err)
	}

	blobCli, err := initBlob(ctx, cfg)
	if err != nil {
		panic(err)
	}

	if cfg.App.Environment == config.APP_ENVIRONMENT_PROD {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}

	router := initGin()
	socket := initSocket()

	lib := &models.Lib{
		Mongo:      mongoCli,
		Blob:       blobCli,
		SocketConn: models.NewSocketConn(),
	}

	return &Server{
		cfg:    cfg,
		lib:    lib,
		router: router,
		socket: socket,
	}
}

func (s *Server) Start(ctx context.Context) {
	log := log.New("server", "Start")

	r := s.router
	s.initCORS()
	s.initSwagger()
	s.initStaticFiles()

	// api v1
	v1 := r.Group("/api/v1")
	s.initHealthCheck(v1)
	s.initHandler(v1)

	// socket
	s.initSocket()

	port := ":" + s.cfg.App.Port
	log.Info("started api successfully with port: %v", port)

	go func() {
		if err := s.router.Run(port); err != nil {
			log.Fatal("router.Run error: %s\n", err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	log := log.New("server", "Stop")

	log.Info("stopping api")
	if err := s.lib.Mongo.Disconnect(ctx); err != nil {
		log.Error("failed to disconnect mongo", err)
	}

	if err := s.socket.Close(); err != nil {
		log.Error("failed to close socket", err)
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

func initSocket() *melody.Melody {
	return melody.New()
}

func initMongo(ctx context.Context, cfg *config.Config) (*mongo.Client, error) {
	log := log.New("server", "initMongo")

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.Mongo.URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Error("mongo.Connect", err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		ip, _ := network.GetPublicIP()
		log.Error(fmt.Sprintf("[%s], client.Ping", ip), err)
		return nil, err
	}
	log.Info("connected to mongo successfully")

	return client, nil
}

func initBlob(_ context.Context, cfg *config.Config) (*models.Blob, error) {
	log := log.New("server", "initBlob")
	host := cfg.AzBlob.Host
	accountKey := cfg.AzBlob.AccountKey
	accountName := cfg.AzBlob.AccountName
	containerName := cfg.Storage.Container
	containerURL := fmt.Sprintf("%s/%s", host, containerName)

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Error("azblob.NewSharedKeyCredential", err)
		return nil, err
	}

	container, err := container.NewClientWithSharedKeyCredential(containerURL, credential, nil)
	if err != nil {
		log.Error("container.NewClientWithSharedKeyCredential", err)
		return nil, err
	}

	service, err := service.NewClientWithSharedKeyCredential(host, credential, nil)
	if err != nil {
		log.Error("service.NewClientWithSharedKeyCredential", err)
		return nil, err
	}

	return &models.Blob{
		Container:  container,
		Credential: credential,
		Service:    service,
	}, nil
}
