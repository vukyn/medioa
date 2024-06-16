package models

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Lib struct {
	Db         *gorm.DB
	Mongo      *mongo.Client
	Blob       *Blob
	SocketConn *SocketConn
}

type Blob struct {
	Container  *container.Client
	Credential *azblob.SharedKeyCredential
	Service    *service.Client
}
