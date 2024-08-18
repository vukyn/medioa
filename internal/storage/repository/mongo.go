package repository

import (
	"context"
	"medioa/config"
	"medioa/constants"
	"medioa/internal/storage/entity"
	commonModel "medioa/models"

	"github.com/vukyn/kuery/conversion"
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
		tableName: (&entity.Storage{}).TableName(),
	}
}

func (m *mongo) withCollection() *mongoo.Collection {
	return m.lib.Mongo.Database(m.cfg.Mongo.Database).Collection(m.tableName)
}

func (m *mongo) GetById(ctx context.Context, id int64) (*entity.Storage, error) {
	return nil, nil
}
func (m *mongo) GetOne(ctx context.Context, queries map[string]interface{}) (*entity.Storage, error) {
	filter := make([]bson.E, 0)
	uuid := conversion.ReadInterfaceV2(queries, constants.FIELD_STORAGE_UUID, "", true)
	downloadUrl := conversion.ReadInterfaceV2(queries, constants.FIELD_STORAGE_DOWNLOAD_URL, "", true)
	downloadPassword := conversion.ReadInterfaceV2(queries, constants.FIELD_STORAGE_DOWNLOAD_PASSWORD, "", true)
	_type := conversion.ReadInterfaceV2(queries, constants.FIELD_STORAGE_TYPE, "", true)
	token := conversion.ReadInterfaceV2(queries, constants.FIELD_STORAGE_TOKEN, "", true)
	ext := conversion.ReadInterfaceV2(queries, constants.FIELD_STORAGE_EXT, "", true)
	secretId := conversion.ReadInterfaceV2(queries, constants.FIELD_STORAGE_SECRET_ID, "", true)
	lifeTime := conversion.ReadInterfaceV2(queries, constants.FIELD_STORAGE_LIFE_TIME, 0, true)

	if uuid != "" {
		filter = append(filter, bson.E{Key: "_id", Value: uuid})
	}
	if downloadUrl != "" {
		filter = append(filter, bson.E{Key: "download_url", Value: downloadUrl})
	}
	if downloadPassword != "" {
		filter = append(filter, bson.E{Key: "download_password", Value: downloadPassword})
	}
	if _type != "" {
		filter = append(filter, bson.E{Key: "type", Value: _type})
	}
	if token != "" {
		filter = append(filter, bson.E{Key: "token", Value: token})
	}
	if lifeTime != 0 {
		filter = append(filter, bson.E{Key: "life_time", Value: lifeTime})
	}
	if ext != "" {
		filter = append(filter, bson.E{Key: "ext", Value: ext})
	}
	if secretId != "" {
		filter = append(filter, bson.E{Key: "secret_id", Value: secretId})
	}

	var obj entity.Storage
	err := m.withCollection().FindOne(ctx, filter).Decode(&obj)
	if err == mongoo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &obj, nil
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
	_, err := m.withCollection().UpdateOne(ctx, bson.D{{Key: "_id", Value: obj.UUID}}, bson.D{{Key: "$set", Value: obj.ToBson()}})
	if err != nil {
		return nil, err
	}
	return obj, nil
}
func (m *mongo) UpdateMany(ctx context.Context, objs []*entity.Storage) (int64, error) {
	return 0, nil
}
