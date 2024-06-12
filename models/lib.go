package models

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Lib struct {
	Db            *gorm.DB
	Mongo         *mongo.Client
	BlobContainer *container.Client
}
