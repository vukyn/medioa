// Đặt tên biến, tên hàm phải sát với ý nghĩa, chức năng.
// Không gán cứng các giá trị số và chuỗi. Nên đặt tên, sử dụng constant để truyền đạt giá trị và ý nghĩa của biến.

package entity

import (
	"medioa/internal/secret/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type Secret struct {
	Id          int64     `gorm:"primarykey;column:id" bson:"id"`
	UUID        string    `gorm:"column:uuid" bson:"_id"`
	Username    string    `gorm:"column:username" bson:"username"`
	Password    string    `gorm:"column:password" bson:"password"`
	PinCode     string    `gorm:"column:pin_code" bson:"pin_code"`
	AccessToken string    `gorm:"column:access_token" bson:"access_token"`
	Type        string    `gorm:"column:type" bson:"type"`
	IsMaster    bool      `gorm:"column:is_master" bson:"is_master"`
	CreatedBy   int64     `gorm:"column:created_by" bson:"created_by"`
	CreatedAt   time.Time `gorm:"autoCreateTime" bson:"created_at"`
}

func (Secret) TableName() string {
	return "secrets"
}

func (e *Secret) Export() *models.Response {
	return &models.Response{
		Id:          e.Id,
		UUID:        e.UUID,
		Username:    e.Username,
		Password:    e.Password,
		PinCode:     e.PinCode,
		AccessToken: e.AccessToken,
		Type:        e.Type,
		IsMaster:    e.IsMaster,
		CreatedBy:   e.CreatedBy,
		CreatedAt:   e.CreatedAt,
	}
}

func (e *Secret) ExportList(objs []*Secret) []*models.Response {
	res := make([]*models.Response, 0)
	for _, obj := range objs {
		res = append(res, obj.Export())
	}
	return res
}

func (e *Secret) ParseFromSaveRequest(req *models.SaveRequest) {
	if req != nil {
		e.Id = req.Id
		e.UUID = req.UUID
		e.Username = req.Username
		e.Password = req.Password
		e.PinCode = req.PinCode
		e.AccessToken = req.AccessToken
		e.Type = req.Type
		e.IsMaster = req.IsMaster
		e.CreatedBy = req.CreatedBy
		e.CreatedAt = req.CreatedAt
	}
}

func (e *Secret) ParseForCreate(req *models.SaveRequest, userId int64) {
	e.ParseFromSaveRequest(req)
	e.hashPassword()
	e.CreatedBy = userId
	e.CreatedAt = time.Now()
}

func (e *Secret) ParseForCreateMany(reqs []*models.SaveRequest, userId int64) []*Secret {
	objs := make([]*Secret, 0)
	for _, v := range reqs {
		obj := &Secret{}
		obj.ParseForCreate(v, userId)
		objs = append(objs, obj)
	}
	return objs
}

func (e *Secret) ParseForUpdate(req *models.SaveRequest, userId int64) {
	e.ParseFromSaveRequest(req)
}

func (e *Secret) ParseForUpdateMany(reqs []*models.SaveRequest, userId int64) []*Secret {
	objs := make([]*Secret, 0)
	for _, v := range reqs {
		obj := &Secret{}
		obj.ParseForUpdate(v, userId)
		objs = append(objs, obj)
	}
	return objs
}

func (e *Secret) ToBson() bson.D {
	d := make(bson.D, 0)
	if e.Id > 0 {
		d = append(d, bson.E{Key: "id", Value: e.Id})
	}
	if e.UUID != "" {
		d = append(d, bson.E{Key: "_id", Value: e.UUID})
	}
	if e.Username != "" {
		d = append(d, bson.E{Key: "username", Value: e.Username})
	}
	if e.Password != "" {
		d = append(d, bson.E{Key: "password", Value: e.Password})
	}
	if e.PinCode != "" {
		d = append(d, bson.E{Key: "pin_code", Value: e.PinCode})
	}
	if e.AccessToken != "" {
		d = append(d, bson.E{Key: "access_token", Value: e.AccessToken})
	}
	if e.Type != "" {
		d = append(d, bson.E{Key: "type", Value: e.Type})
	}
	if e.IsMaster {
		d = append(d, bson.E{Key: "is_master", Value: true})
	}
	if e.CreatedBy > 0 {
		d = append(d, bson.E{Key: "created_by", Value: e.CreatedBy})
	}
	if !e.CreatedAt.IsZero() {
		d = append(d, bson.E{Key: "created_at", Value: e.CreatedAt.UnixMilli()})
	}
	return d
}

func (s *Secret) hashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	s.Password = string(hashedPassword)
	return nil
}
