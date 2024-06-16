package models

import "mime/multipart"

type UploadRequest struct {
	SessionId string                `json:"id"`
	File      *multipart.FileHeader `json:"file"`
}

func (r *UploadRequest) ToBlobRequest() *UploadBlobRequest {
	return &UploadBlobRequest{
		SessionId: r.SessionId,
		File:      r.File,
	}
}

type UploadResponse struct {
	Url      string `json:"url"`
	Token    string `json:"token"`
	Ext      string `json:"ext"`
	FileName string `json:"file_name"` // review to remove later
}

type DownloadRequest struct {
	FileName string `json:"file_name"`
	Token    string `json:"token"`
}

type DownloadResponse struct {
	Url string `json:"url"`
}
