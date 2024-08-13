package config

import (
	"fmt"
	"medioa/pkg/log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/vukyn/kuery/crypto"
)

const (
	APP_ENVIRONMENT_LOCAL = "local"
	APP_ENVIRONMENT_PROD  = "prod"
)

type Config struct {
	App      AppConfig
	Log      log.Config
	Mongo    MongoConfig
	AzAd     AzAdConfig
	AzBlob   AzBlobConfig
	Storage  StorageConfig
	Cors     CorsConfig
	Secret   SecretConfig
	Upload   UploadConfig
	Download DownloadConfig
}

type AppConfig struct {
	Version         string
	Environment     string
	Port            string
	Host            string
	ShutdownTimeout int
}

type CorsConfig struct {
	AllowOrigins []string
	AllowHeaders []string
	AllowMethods []string
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

type AzAdConfig struct {
	TenantId     string
	ClientId     string
	ClientSecret string
}

type StorageConfig struct {
	Container string
}

type SecretConfig struct {
	SecretKey string
}

type UploadConfig struct {
	MaxSizeMB int64
}

type DownloadConfig struct {
	Expire int64 // in days
}

func Load() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	cfg := &Config{}
	parseAppConfig(cfg)
	parseLogConfig(cfg)
	parseMongoConfig(cfg)
	parseAzBlobConfig(cfg)
	parseStorageConfig(cfg)
	parseCorsConfig(cfg)
	parseAzAdConfig(cfg)
	parseSecretConfig(cfg)
	parseUploadConfig(cfg)
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

func parseCorsConfig(cfg *Config) {
	cfg.Cors.AllowOrigins = strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ",")
	cfg.Cors.AllowHeaders = strings.Split(os.Getenv("CORS_ALLOW_HEADERS"), ",")
	cfg.Cors.AllowMethods = strings.Split(os.Getenv("CORS_ALLOW_METHODS"), ",")
}

func parseAzAdConfig(cfg *Config) {
	cfg.AzAd.TenantId = os.Getenv("AZAD_TENANT_ID")
	cfg.AzAd.ClientId = os.Getenv("AZAD_CLIENT_ID")
	cfg.AzAd.ClientSecret = os.Getenv("AZAD_CLIENT_SECRET")
}

func parseSecretConfig(cfg *Config) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey != "" {
		cfg.Secret.SecretKey = string(crypto.Md5Encrypt([]byte(secretKey)))
	}
}

func parseUploadConfig(cfg *Config) {
	maxSizeMB, _ := strconv.ParseInt(os.Getenv("UPLOAD_MAX_SIZE_MB"), 10, 64)
	cfg.Upload.MaxSizeMB = maxSizeMB
}

func parseDownloadConfig(cfg *Config) {
	expire, _ := strconv.ParseInt(os.Getenv("DOWNLOAD_EXPIRE"), 10, 64)
	cfg.Download.Expire = expire
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

	// if len(cfg.Cors.AllowOrigins) == 0 {
	// 	return fmt.Errorf("cors allow origins is required")
	// }

	// if len(cfg.Cors.AllowHeaders) == 0 {
	// 	return fmt.Errorf("cors allow headers is required")
	// }

	// if len(cfg.Cors.AllowMethods) == 0 {
	// 	return fmt.Errorf("cors allow methods is required")
	// }

	// if cfg.AzAd.TenantId == "" {
	// 	return fmt.Errorf("azad tenant id is required")
	// }

	// if cfg.AzAd.ClientId == "" {
	// 	return fmt.Errorf("azad client id is required")
	// }

	// if cfg.AzAd.ClientSecret == "" {
	// 	return fmt.Errorf("azad client secret is required")
	// }

	if cfg.Secret.SecretKey == "" {
		return fmt.Errorf("secret key is required")
	}

	if cfg.Upload.MaxSizeMB <= 0 {
		return fmt.Errorf("upload max size mb is invalid")
	}

	if cfg.Download.Expire <= 0 {
		return fmt.Errorf("download expire is invalid")
	}

	return nil
}
