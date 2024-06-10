package models

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Lib struct {
	Db    *gorm.DB
	Mongo *mongo.Client
}
