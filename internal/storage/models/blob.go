package models

import "mime/multipart"

type UploadBlobRequest struct {
	SessionId string
	File      *multipart.FileHeader
}

type UploadBlobResponse struct {
	Url      string
	Token    string
	Ext      string
	FileName string
}

type DownloadSASRequest struct {
	FileName string
}

type DownloadSASResponse struct {
	Url string
}
