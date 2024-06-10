package repository

import (
	"context"
	"medioa/config"
	"medioa/internal/storage/entity"
	commonModel "medioa/models"

	mongoo "go.mongodb.org/mongo-driver/mongo"
)

type mongo struct {
	cfg       *config.Config
	lib       *commonModel.Lib
	tableName string
}

func InitMongo(cfg *config.Config, lib *commonModel.Lib) IRepository {
	return &mongo{
		cfg:       cfg,
		lib:       lib,
		tableName: (&entity.Storage{}).TableName(),
	}
}

func (m *mongo) withCollection() *mongoo.Collection {
	return m.lib.Mongo.Database(m.cfg.MongoConfig.Database).Collection(m.tableName)
}

func (m *mongo) GetById(ctx context.Context, id int64) (*entity.Storage, error) {
	return nil, nil
}
func (m *mongo) GetOne(ctx context.Context, queries map[string]interface{}) (*entity.Storage, error) {
	return nil, nil
}
func (m *mongo) GetList(ctx context.Context, queries map[string]interface{}) ([]*entity.Storage, error) {
	return nil, nil
}
func (m *mongo) GetListPaging(ctx context.Context, queries map[string]interface{}) ([]*entity.Storage, error) {
	return nil, nil
}
func (m *mongo) Count(ctx context.Context, queries map[string]interface{}) (int64, error) {
	return 0, nil
}
func (m *mongo) Create(ctx context.Context, obj *entity.Storage) (*entity.Storage, error) {
	_, err := m.withCollection().InsertOne(ctx, obj.ToBson())
	if err != nil {
		return nil, err
	}
	return obj, nil
}
func (m *mongo) CreateMany(ctx context.Context, objs []*entity.Storage) ([]*entity.Storage, error) {
	return nil, nil
}
func (m *mongo) Update(ctx context.Context, obj *entity.Storage) (*entity.Storage, error) {
	return nil, nil
}
func (m *mongo) UpdateMany(ctx context.Context, objs []*entity.Storage) (int64, error) {
	return 0, nil
}
