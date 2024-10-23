package repository

import (
	"context"
	"medioa/config"
	"medioa/constants"
	"medioa/internal/secret/entity"
	commonModel "medioa/models"

	"github.com/vukyn/kuery/conv"
	"go.mongodb.org/mongo-driver/bson"
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
		tableName: (&entity.Secret{}).TableName(),
	}
}

func (m *mongo) withCollection() *mongoo.Collection {
	return m.lib.Mongo.Database(m.cfg.Mongo.Database).Collection(m.tableName)
}

func (m *mongo) GetById(ctx context.Context, id int64) (*entity.Secret, error) {
	return nil, nil
}
func (m *mongo) GetOne(ctx context.Context, queries map[string]any) (*entity.Secret, error) {
	filter := make([]bson.E, 0)
	uuid := conv.ReadInterface(queries, constants.FIELD_STORAGE_UUID, "")
	username := conv.ReadInterface(queries, constants.FIELD_SECRET_USERNAME, "")
	accessToken := conv.ReadInterface(queries, constants.FIELD_SECRET_ACCESS_TOKEN, "")
	_type := conv.ReadInterface(queries, constants.FIELD_SECRET_TYPE, "")

	if uuid != "" {
		filter = append(filter, bson.E{Key: "_id", Value: uuid})
	}
	if username != "" {
		filter = append(filter, bson.E{Key: "username", Value: username})
	}
	if accessToken != "" {
		filter = append(filter, bson.E{Key: "access_token", Value: accessToken})
	}
	if _type != "" {
		filter = append(filter, bson.E{Key: "type", Value: _type})
	}

	var obj entity.Secret
	err := m.withCollection().FindOne(ctx, filter).Decode(&obj)
	if err == mongoo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &obj, nil
}
func (m *mongo) GetList(ctx context.Context, queries map[string]any) ([]*entity.Secret, error) {
	return nil, nil
}
func (m *mongo) GetListPaging(ctx context.Context, queries map[string]any) ([]*entity.Secret, error) {
	return nil, nil
}
func (m *mongo) Count(ctx context.Context, queries map[string]any) (int64, error) {
	return 0, nil
}
func (m *mongo) Create(ctx context.Context, obj *entity.Secret) (*entity.Secret, error) {
	_, err := m.withCollection().InsertOne(ctx, obj.ToBson())
	if err != nil {
		return nil, err
	}
	return obj, nil
}
func (m *mongo) CreateMany(ctx context.Context, objs []*entity.Secret) ([]*entity.Secret, error) {
	return nil, nil
}
func (m *mongo) Update(ctx context.Context, obj *entity.Secret) (*entity.Secret, error) {
	_, err := m.withCollection().UpdateOne(ctx, bson.D{{Key: "_id", Value: obj.UUID}}, bson.D{{Key: "$set", Value: obj.ToBson()}})
	if err != nil {
		return nil, err
	}
	return obj, nil
}
func (m *mongo) UpdateMany(ctx context.Context, objs []*entity.Secret) (int64, error) {
	return 0, nil
}
