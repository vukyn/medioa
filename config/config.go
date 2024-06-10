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
	AppConfig   AppConfig
	LogConfig   log.Config
	MongoConfig MongoConfig
}

type AppConfig struct {
	Version         string
	Environment     string
	Port            string
	ShutdownTimeout int
}

type MongoConfig struct {
	URI      string
	Database string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	parseAppConfig(cfg)
	parseLogConfig(cfg)
	parseMongoConfig(cfg)
	validation(cfg)

	return cfg, nil
}

func parseLogConfig(cfg *Config) {
	cfg.LogConfig.Mode = os.Getenv("LOG_MODE")
	cfg.LogConfig.Level = os.Getenv("LOG_LEVEL")
}

func parseAppConfig(cfg *Config) {
	cfg.AppConfig.Version = os.Getenv("VERSION")
	cfg.AppConfig.Environment = os.Getenv("ENVIRONMENT")
	cfg.AppConfig.Port = os.Getenv("PORT")
	shutdownTimeout, _ := strconv.Atoi(os.Getenv("SHUTDOWN_TIMEOUT"))
	cfg.AppConfig.ShutdownTimeout = shutdownTimeout
}

func parseMongoConfig(cfg *Config) {
	cfg.MongoConfig.URI = os.Getenv("MONGO_URI")
	cfg.MongoConfig.Database = os.Getenv("MONGO_DATABASE")
}

func validation(cfg *Config) error {
	if cfg.AppConfig.Version == "" {
		return fmt.Errorf("version is required")
	}

	if cfg.AppConfig.Environment == "" {
		return fmt.Errorf("environment is required")
	}
	switch cfg.AppConfig.Environment {
	case APP_ENVIRONMENT_LOCAL, APP_ENVIRONMENT_PROD:
	default:
		return fmt.Errorf("environment is invalid")
	}

	if cfg.AppConfig.Port == "" {
		return fmt.Errorf("port is required")
	}

	if cfg.AppConfig.ShutdownTimeout < 0 {
		return fmt.Errorf("shutdownTimeout is invalid")
	}

	if cfg.MongoConfig.URI == "" {
		return fmt.Errorf("mongo URI is required")
	}

	if cfg.MongoConfig.Database == "" {
		return fmt.Errorf("mongo database is required")
	}

	return nil
}
