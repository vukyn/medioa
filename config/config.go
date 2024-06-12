package config

import (
	"fmt"
	"medioa/pkg/log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	APP_ENVIRONMENT_LOCAL = "local"
	APP_ENVIRONMENT_PROD  = "prod"
)

type Config struct {
	App     AppConfig
	Log     log.Config
	Mongo   MongoConfig
	AzBlob  AzBlobConfig
	Storage StorageConfig
}

type AppConfig struct {
	Version         string
	Environment     string
	Port            string
	Host            string
	ShutdownTimeout int
}

type MongoConfig struct {
	URI      string
	Database string
}

type AzBlobConfig struct {
	Host        string
	AccountName string
	AccountKey  string
}

type StorageConfig struct {
	Container string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	parseAppConfig(cfg)
	parseLogConfig(cfg)
	parseMongoConfig(cfg)
	parseAzBlobConfig(cfg)
	parseStorageConfig(cfg)
	validation(cfg)

	return cfg, nil
}

func parseLogConfig(cfg *Config) {
	cfg.Log.Mode = os.Getenv("LOG_MODE")
	cfg.Log.Level = os.Getenv("LOG_LEVEL")
}

func parseAppConfig(cfg *Config) {
	cfg.App.Version = os.Getenv("VERSION")
	cfg.App.Environment = os.Getenv("ENVIRONMENT")
	cfg.App.Port = os.Getenv("PORT")
	cfg.App.Host = os.Getenv("HOST")
	shutdownTimeout, _ := strconv.Atoi(os.Getenv("SHUTDOWN_TIMEOUT"))
	cfg.App.ShutdownTimeout = shutdownTimeout
}

func parseMongoConfig(cfg *Config) {
	cfg.Mongo.URI = os.Getenv("MONGO_URI")
	cfg.Mongo.Database = os.Getenv("MONGO_DATABASE")
}

func parseAzBlobConfig(cfg *Config) {
	cfg.AzBlob.Host = os.Getenv("AZBLOB_HOST")
	cfg.AzBlob.AccountName = os.Getenv("AZBLOB_ACCOUNT_NAME")
	cfg.AzBlob.AccountKey = os.Getenv("AZBLOB_ACCOUNT_KEY")
}

func parseStorageConfig(cfg *Config) {
	cfg.Storage.Container = os.Getenv("STORAGE_CONTAINER")
}

func validation(cfg *Config) error {
	if cfg.App.Version == "" {
		return fmt.Errorf("version is required")
	}

	if cfg.App.Environment == "" {
		return fmt.Errorf("environment is required")
	}
	switch cfg.App.Environment {
	case APP_ENVIRONMENT_LOCAL, APP_ENVIRONMENT_PROD:
	default:
		return fmt.Errorf("environment is invalid")
	}

	if cfg.App.Port == "" {
		return fmt.Errorf("port is required")
	}

	if cfg.App.Host == "" {
		return fmt.Errorf("host is required")
	}

	if cfg.App.ShutdownTimeout < 0 {
		return fmt.Errorf("shutdownTimeout is invalid")
	}

	if cfg.Mongo.URI == "" {
		return fmt.Errorf("mongo URI is required")
	}

	if cfg.Mongo.Database == "" {
		return fmt.Errorf("mongo database is required")
	}

	if cfg.AzBlob.Host == "" {
		return fmt.Errorf("azblob host is required")
	}

	if cfg.AzBlob.AccountName == "" {
		return fmt.Errorf("azblob account name is required")
	}

	if cfg.AzBlob.AccountKey == "" {
		return fmt.Errorf("azblob account key is required")
	}

	if cfg.Storage.Container == "" {
		return fmt.Errorf("storage container is required")
	}

	return nil
}
