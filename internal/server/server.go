package server

import (
	"context"
	"fmt"
	"io"
	"medioa/config"
	"medioa/models"
	"medioa/pkg/log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "medioa/docs"
)

type Server struct {
	lib *models.Lib
	cfg *config.Config
}

func New(ctx context.Context, cfg *config.Config) *Server {
	mongoCli, err := initMongo(ctx, cfg)
	if err != nil {
		panic(err)
	}

	blobContainerCli, err := initBlobContainer(ctx, cfg)
	if err != nil {
		panic(err)
	}

	lib := &models.Lib{
		Mongo:         mongoCli,
		BlobContainer: blobContainerCli,
	}
	return &Server{
		cfg: cfg,
		lib: lib,
	}
}

func (s *Server) Start(ctx context.Context) {
	log := log.New("server", "Start")

	port := s.cfg.App.Port
	appEnv := s.cfg.App.Environment

	if appEnv == config.APP_ENVIRONMENT_PROD {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}

	r := initGin()
	s.initSwagger(r)

	// api v1
	v1 := r.Group("/api/v1")
	s.initHealthCheck(v1)
	s.initHandler(v1)

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
		log.Error("client.Ping", err)
		return nil, err
	}
	log.Info("connected to mongo successfully")

	return client, nil
}

func initBlobContainer(_ context.Context, cfg *config.Config) (*container.Client, error) {
	log := log.New("server", "initBlobContainer")
	host := cfg.AzBlob.Host
	accountKey := cfg.AzBlob.AccountKey
	accountName := cfg.AzBlob.AccountName
	containerName := cfg.Storage.Container

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Error("azblob.NewSharedKeyCredential", err)
		return nil, err
	}

	containerURL := fmt.Sprintf("%s/%s", host, containerName)
	client, err := container.NewClientWithSharedKeyCredential(containerURL, credential, nil)
	if err != nil {
		log.Error("container.NewClientWithSharedKeyCredential", err)
		return nil, err
	}

	return client, nil
}
