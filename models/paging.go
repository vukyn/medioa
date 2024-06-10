package models

type ListPaging struct {
	Page    int64
	Size    int64
	Count   int64
	Records any
}

type RequestParams struct {
	Page         int64
	Size         int64
	SortBy       string
	OrderBy      string
	SortMultiple string
}
