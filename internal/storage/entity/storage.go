// Đặt tên biến, tên hàm phải sát với ý nghĩa, chức năng.
// Không gán cứng các giá trị số và chuỗi. Nên đặt tên, sử dụng constant để truyền đạt giá trị và ý nghĩa của biến.

package entity

import (
	"medioa/internal/storage/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Storage struct {
	Id          int64     `gorm:"primarykey;column:id" bson:"id"`
	UUID        string    `gorm:"column:uuid" bson:"_id"`
	DownloadUrl string    `gorm:"column:download_url" bson:"download_url"`
	Type        string    `gorm:"column:type" bson:"type"`
	Token       string    `gorm:"column:token;default:(-)" bson:"token"`
	LifeTime    int       `gorm:"column:life_time;default:(-)" bson:"life_time"`
	FileName    string    `gorm:"column:file_name" bson:"file_name"`
	FileSize    int64     `gorm:"column:file_size" bson:"file_size"`
	Ext         string    `gorm:"column:ext" bson:"ext"`
	SecretId    string    `gorm:"column:secret_id" bson:"secret_id"`
	CreatedBy   int64     `gorm:"column:created_by" bson:"created_by"`
	CreatedAt   time.Time `gorm:"autoCreateTime" bson:"created_at"`
}

func (s *Storage) TableName() string {
	return "storage"
}

func (e *Storage) Export() *models.Response {
	return &models.Response{
		Id:          e.Id,
		UUID:        e.UUID,
		DownloadUrl: e.DownloadUrl,
		Type:        e.Type,
		Token:       e.Token,
		LifeTime:    e.LifeTime,
		FileName:    e.FileName,
		FileSize:    e.FileSize,
		Ext:         e.Ext,
		SecretId:    e.SecretId,
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
		e.FileName = req.FileName
		e.FileSize = req.FileSize
		e.Ext = req.Ext
		e.SecretId = req.SecretId
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
	d := make(bson.D, 0)
	if e.Id > 0 {
		d = append(d, bson.E{Key: "id", Value: e.Id})
	}
	if e.UUID != "" {
		d = append(d, bson.E{Key: "_id", Value: e.UUID})
	}
	if e.DownloadUrl != "" {
		d = append(d, bson.E{Key: "download_url", Value: e.DownloadUrl})
	}
	if e.Type != "" {
		d = append(d, bson.E{Key: "type", Value: e.Type})
	}
	if e.Token != "" {
		d = append(d, bson.E{Key: "token", Value: e.Token})
	}
	if e.LifeTime > 0 {
		d = append(d, bson.E{Key: "life_time", Value: e.LifeTime})
	}
	if e.FileName != "" {
		d = append(d, bson.E{Key: "file_name", Value: e.FileName})
	}
	if e.FileSize > 0 {
		d = append(d, bson.E{Key: "file_size", Value: e.FileSize})
	}
	if e.Ext != "" {
		d = append(d, bson.E{Key: "ext", Value: e.Ext})
	}
	if e.SecretId != "" {
		d = append(d, bson.E{Key: "secret_id", Value: e.SecretId})
	}
	if e.CreatedBy > 0 {
		d = append(d, bson.E{Key: "created_by", Value: e.CreatedBy})
	}
	if !e.CreatedAt.IsZero() {
		d = append(d, bson.E{Key: "created_at", Value: e.CreatedAt.UnixMilli()})
	}
	return d
}
