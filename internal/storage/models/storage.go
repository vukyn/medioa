package models

import "mime/multipart"

type UploadRequest struct {
	SessionId string                `json:"id"`
	File      *multipart.FileHeader `json:"file"`
	FileName  string                `json:"filename"`
}

func (r *UploadRequest) ToBlobRequest() *UploadBlobRequest {
	return &UploadBlobRequest{
		SessionId: r.SessionId,
		File:      r.File,
	}
}

type UploadWithSecretRequest struct {
	SessionId string                `json:"id"`
	Secret    string                `json:"secret"`
	File      *multipart.FileHeader `json:"file"`
	FileName  string                `json:"filename"`
}

func (r *UploadWithSecretRequest) ToBlobRequest() *UploadBlobRequest {
	return &UploadBlobRequest{
		SessionId: r.SessionId,
		File:      r.File,
	}
}

type UploadResponse struct {
	Url      string `json:"url"`
	Token    string `json:"token"`
	Ext      string `json:"ext"`
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}

type DownloadRequest struct {
	FileId string `json:"file_id"`
	Token  string `json:"token"`
}

type DownloadWithSecretRequest struct {
	FileId string `json:"file_id"`
	Secret string `json:"secret"`
	Token  string `json:"token"`
}

type DownloadResponse struct {
	Url string `json:"url"`
}
