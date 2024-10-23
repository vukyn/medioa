// Đặt tên biến, tên hàm phải sát với ý nghĩa, chức năng.
// Không gán cứng các giá trị số và chuỗi. Nên đặt tên, sử dụng constant để truyền đạt giá trị và ý nghĩa của biến.

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
	Id          int64
	UUID        string
	Username    string
	Password    string
	PinCode     string
	AccessToken string
	Type        string
	CreatedBy   int64
}

func (r *RequestParams) trimSpace() {
	r.UUID = strings.TrimSpace(r.UUID)
	r.Username = strings.TrimSpace(r.Username)
	r.Password = strings.TrimSpace(r.Password)
	r.PinCode = strings.TrimSpace(r.PinCode)
	r.AccessToken = strings.TrimSpace(r.AccessToken)
	r.Type = strings.TrimSpace(r.Type)
}
func (r *RequestParams) ToMap() map[string]any {
	r.trimSpace()

	if strings.ToLower(r.OrderBy) != constants.SORT_ORDER_ASC {
		r.OrderBy = constants.SORT_ORDER_DESC
	}
	return map[string]any{
		constants.FIELD_SECRET_ID:           r.Id,
		constants.FIELD_SECRET_UUID:         r.UUID,
		constants.FIELD_SECRET_USERNAME:     r.Username,
		constants.FIELD_SECRET_PASSWORD:     r.Password,
		constants.FIELD_SECRET_PIN_CODE:     r.PinCode,
		constants.FIELD_SECRET_ACCESS_TOKEN: r.AccessToken,
		constants.FIELD_SECRET_TYPE:         r.Type,
		constants.FIELD_SECRET_CREATED_BY:   r.CreatedBy,
		constants.FIELD_PAGE:                r.Page,
		constants.FIELD_SIZE:                r.Size,
		constants.FIELD_ORDER_BY:            r.OrderBy,
		constants.FIELD_SORT_BY:             r.SortBy,
		constants.FIELD_SORT_MULTIPLE:       r.SortMultiple,
	}
}

type Response struct {
	Id          int64
	UUID        string
	Username    string
	Password    string
	PinCode     string
	AccessToken string
	Type        string
	IsMaster    bool
	CreatedBy   int64
	CreatedAt   time.Time
}

type SaveRequest struct {
	Id          int64
	UUID        string
	Username    string
	Password    string
	PinCode     string
	AccessToken string
	Type        string
	IsMaster    bool
	CreatedBy   int64
	CreatedAt   time.Time
}

type ListPaging struct {
	commonModel.ListPaging
	Records []*Response
}
