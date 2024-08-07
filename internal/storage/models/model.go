// Đặt tên biến, tên hàm phải sát với ý nghĩa, chức năng.
// Không gán cứng các giá trị số và chuỗi. Nên đặt tên, sử dụng constants để truyền đạt giá trị và ý nghĩa của biến.

package models

import (
	"medioa/constants"
	commonModel "medioa/models"
	"strings"
	"time"
)

type RequestParams struct {
	commonModel.RequestParams
	ConfigQuery int
	Id          int
	UUID        string
	DownloadUrl string
	Type        string
	Token       string
	Ext         string
	LifeTime    int64
	SecretId    string
	CreatedBy   int64
}

func (r *RequestParams) trimSpace() {
	r.DownloadUrl = strings.TrimSpace(r.DownloadUrl)
	r.Type = strings.TrimSpace(r.Type)
	r.Token = strings.TrimSpace(r.Token)
	r.Ext = strings.TrimSpace(r.Ext)
	r.SecretId = strings.TrimSpace(r.SecretId)
}
func (r *RequestParams) ToMap() map[string]interface{} {
	r.trimSpace()

	if strings.ToLower(r.OrderBy) != constants.SORT_ORDER_ASC {
		r.OrderBy = constants.SORT_ORDER_DESC
	}
	return map[string]interface{}{
		constants.FIELD_STORAGE_ID:           r.Id,
		constants.FIELD_STORAGE_UUID:         r.UUID,
		constants.FIELD_STORAGE_DOWNLOAD_URL: r.DownloadUrl,
		constants.FIELD_STORAGE_TYPE:         r.Type,
		constants.FIELD_STORAGE_TOKEN:        r.Token,
		constants.FIELD_STORAGE_LIFE_TIME:    r.LifeTime,
		constants.FIELD_STORAGE_EXT:          r.Ext,
		constants.FIELD_STORAGE_SECRET_ID:    r.SecretId,
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
	UUID        string    `json:"uuid"`
	DownloadUrl string    `json:"download_url"`
	Type        string    `json:"type"`
	Token       string    `json:"token"`
	LifeTime    int64     `json:"life_time"`
	FileName    string    `json:"file_name"`
	FileSize    int64     `json:"file_size"`
	Ext         string    `json:"ext"`
	SecretId    string    `json:"secret_id"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	ChunkIds    []string  `json:"chunk_ids"`
	TotalChunks int64     `json:"total_chunks"`
}

type SaveRequest struct {
	Id          int64
	UUID        string
	DownloadUrl string
	Type        string
	Token       string
	FileName    string
	FileSize    int64
	Ext         string
	SecretId    string
	LifeTime    int64
	ChunkIds    *[]string
	TotalChunks int64
	CreatedBy   int64
	CreatedAt   time.Time
}

type ListPaging struct {
	commonModel.ListPaging
	Records []*Response
}

type AddChunkRequest struct {
	Id      string
	ChunkId string
}
