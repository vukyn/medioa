// Đặt tên biến, tên hàm phải sát với ý nghĩa, chức năng.
// Không gán cứng các giá trị số và chuỗi. Nên đặt tên, sử dụng constants để truyền đạt giá trị và ý nghĩa của biến.

package models

import (
	"medioa/constants"
	commonModel "medioa/models"
	"mime/multipart"
	"strings"
	"time"
)

type RequestParams struct {
	commonModel.RequestParams
	ConfigQuery int
	Id          int
	DownloadUrl string
	Type        string
	Token       string
	LifeTime    int
	CreatedBy   int
}

func (r *RequestParams) trimSpace() {
	r.DownloadUrl = strings.TrimSpace(r.DownloadUrl)
	r.Type = strings.TrimSpace(r.Type)
	r.Token = strings.TrimSpace(r.Token)
}
func (r *RequestParams) ToMap() map[string]interface{} {
	r.trimSpace()

	if strings.ToLower(r.OrderBy) != constants.SORT_ORDER_ASC {
		r.OrderBy = constants.SORT_ORDER_DESC
	}
	return map[string]interface{}{
		constants.FIELD_STORAGE_ID:           r.Id,
		constants.FIELD_STORAGE_DOWNLOAD_URL: r.DownloadUrl,
		constants.FIELD_STORAGE_TYPE:         r.Type,
		constants.FIELD_STORAGE_TOKEN:        r.Token,
		constants.FIELD_STORAGE_LIFE_TIME:    r.LifeTime,
		constants.FIELD_STORAGE_CREATED_BY:   r.CreatedBy,
		constants.FIELD_PAGE:                 r.Page,
		constants.FIELD_SIZE:                 r.Size,
		constants.FIELD_ORDER_BY:             r.OrderBy,
		constants.FIELD_SORT_BY:              r.SortBy,
		constants.FIELD_SORT_MULTIPLE:        r.SortMultiple,
	}
}

type Response struct {
	Id          int64     `json:"id"`
	DownloadUrl string    `json:"download_url"`
	Type        string    `json:"type"`
	Token       string    `json:"token"`
	LifeTime    int       `json:"life_time"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type SaveRequest struct {
	Id          int64
	DownloadUrl string
	Type        string
	Token       string
	LifeTime    int
	CreatedBy   int
	CreatedAt   time.Time
}

type ListPaging struct {
	commonModel.ListPaging
	Records []*Response
}

type UploadRequest struct {
	File *multipart.FileHeader
}
