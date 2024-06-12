package models

import "mime/multipart"

type UploadRequest struct {
	File *multipart.FileHeader `json:"file"`
}

type UploadResponse struct {
	Url      string `json:"url"`
	Token    string `json:"token"`
	Ext      string `json:"ext"`
	FileName string `json:"file_name"`
}
