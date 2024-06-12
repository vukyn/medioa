// Đặt tên biến, tên hàm phải sát với ý nghĩa, chức năng.
// Không gán cứng các giá trị số và chuỗi. Nên đặt tên, sử dụng constant để truyền đạt giá trị và ý nghĩa của biến.

package entity

import (
	"medioa/internal/storage/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Storage struct {
	Id          int64     `gorm:"primarykey;column:id"`
	UUID        uuid.UUID `gorm:"column:uuid"`
	DownloadUrl string    `gorm:"column:download_url"`
	Type        string    `gorm:"column:type"`
	Token       string    `gorm:"column:token;default:(-)"`
	LifeTime    int       `gorm:"column:life_time;default:(-)"`
	CreatedBy   int64     `gorm:"column:created_by"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (s *Storage) TableName() string {
	return "storage"
}

func (e *Storage) Export() *models.Response {
	return &models.Response{
		Id:          e.Id,
		DownloadUrl: e.DownloadUrl,
		Type:        e.Type,
		Token:       e.Token,
		LifeTime:    e.LifeTime,
		CreatedBy:   e.CreatedBy,
		CreatedAt:   e.CreatedAt,
	}
}

func (e *Storage) ExportList(objs []*Storage) []*models.Response {
	res := make([]*models.Response, 0)
	for _, obj := range objs {
		res = append(res, obj.Export())
	}
	return res
}

func (e *Storage) ParseFromSaveRequest(req *models.SaveRequest) {
	if req != nil {
		e.Id = req.Id
		e.UUID = req.UUID
		e.DownloadUrl = req.DownloadUrl
		e.Type = req.Type
		e.Token = req.Token
		e.LifeTime = req.LifeTime
		e.CreatedBy = req.CreatedBy
		e.CreatedAt = req.CreatedAt
	}
}

func (e *Storage) ParseForCreate(req *models.SaveRequest, userId int64) {
	e.ParseFromSaveRequest(req)
	e.CreatedBy = userId
	e.CreatedAt = time.Now()
}

func (e *Storage) ParseForCreateMany(reqs []*models.SaveRequest, userId int64) []*Storage {
	objs := make([]*Storage, 0)
	for _, v := range reqs {
		obj := &Storage{}
		obj.ParseForCreate(v, userId)
		objs = append(objs, obj)
	}
	return objs
}

func (e *Storage) ParseForUpdate(req *models.SaveRequest, userId int64) {
	e.ParseFromSaveRequest(req)
}

func (e *Storage) ParseForUpdateMany(reqs []*models.SaveRequest, userId int64) []*Storage {
	objs := make([]*Storage, 0)
	for _, v := range reqs {
		obj := &Storage{}
		obj.ParseForUpdate(v, userId)
		objs = append(objs, obj)
	}
	return objs
}

func (e *Storage) ToBson() bson.D {
	return bson.D{
		{Key: "id", Value: e.Id},
		{Key: "_id", Value: e.UUID.String()},
		{Key: "download_url", Value: e.DownloadUrl},
		{Key: "type", Value: e.Type},
		{Key: "token", Value: e.Token},
		{Key: "life_time", Value: e.LifeTime},
		{Key: "created_by", Value: e.CreatedBy},
		{Key: "created_at", Value: e.CreatedAt.UnixMilli()},
	}
}
