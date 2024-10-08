package models

import (
	"errors"
	azBlobModel "medioa/internal/azblob/models"
	"medioa/pkg/xtype"
)

type GetFileInfoRequest struct {
	FileId string `json:"file_id"`
}

type GetFileInfoResponse struct {
	FileId    string `json:"file_id"`
	FileName  string `json:"file_name"`
	FileSize  int64  `json:"file_size"`
	HasSecret bool   `json:"has_secret"`
}

type UploadRequest struct {
	SessionId string     `json:"session_id"`
	File      xtype.File `json:"file"`
	URL       string     `json:"url"`
	FileName  string     `json:"file_name"`
}

func (r *UploadRequest) Validate() error {
	if r.File == nil && r.URL == "" {
		return errors.New("file or url is required")
	}
	return nil
}

func (r *UploadRequest) ToURLRequest() *azBlobModel.UploadURLRequest {
	return &azBlobModel.UploadURLRequest{
		URL: r.URL,
	}
}

func (r *UploadRequest) ToBlobRequest() *azBlobModel.UploadBlobRequest {
	return &azBlobModel.UploadBlobRequest{
		SessionId: r.SessionId,
		File:      r.File,
	}
}

type UploadChunkRequest struct {
	SessionId   string     `json:"session_id"`
	FileId      string     `json:"file_id"`
	FileName    string     `json:"file_name"`
	Chunk       xtype.File `json:"chunk"`
	ChunkIndex  int64      `json:"chunk_index"`
	TotalChunks int64      `json:"total_chunks"`
}

func (r *UploadChunkRequest) ToBlobRequest(token string) *azBlobModel.UploadChunkRequest {
	return &azBlobModel.UploadChunkRequest{
		SessionId:   r.SessionId,
		Token:       token,
		FileName:    r.FileName,
		Chunk:       r.Chunk,
		ChunkIndex:  r.ChunkIndex,
		TotalChunks: r.TotalChunks,
	}
}

type UploadChunkResponse struct {
	ChunkId string `json:"chunk_id"`
	FileId  string `json:"file_id"`
}

type CommitChunkRequest struct {
	SessionId string `json:"session_id" swaggerignore:"true"`
	Secret    string `json:"secret" swaggerignore:"true"`
	FileId    string `json:"file_id"`
}

type CommitChunkResponse struct {
	Url      string `json:"url"`
	Token    string `json:"token"`
	Ext      string `json:"ext"`
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}

type UploadWithSecretRequest struct {
	SessionId string
	Secret    string
	File      xtype.File
	FileName  string
}

func (r *UploadWithSecretRequest) ToBlobRequest(secretId string) *azBlobModel.UploadBlobRequest {
	return &azBlobModel.UploadBlobRequest{
		SessionId: r.SessionId,
		SecretId:  secretId,
		File:      r.File,
	}
}

type UploadChunkWithSecretRequest struct {
	SessionId   string     `json:"session_id"`
	Secret      string     `json:"secret"`
	FileId      string     `json:"file_id"`
	FileName    string     `json:"file_name"`
	Chunk       xtype.File `json:"chunk"`
	ChunkIndex  int64      `json:"chunk_index"`
	TotalChunks int64      `json:"total_chunks"`
}

func (r *UploadChunkWithSecretRequest) ToBlobRequest(secretId, token string) *azBlobModel.UploadChunkRequest {
	return &azBlobModel.UploadChunkRequest{
		SessionId:   r.SessionId,
		SecretId:    secretId,
		Token:       token,
		FileName:    r.FileName,
		Chunk:       r.Chunk,
		ChunkIndex:  r.ChunkIndex,
		TotalChunks: r.TotalChunks,
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
	FileId           string `json:"file_id"`
	Token            string `json:"token"`
	Secret           string `json:"secret"`
	DownloadPassword string `json:"download_password"`
}

type RequestDownloadRequest struct {
	FileId string `json:"file_id"`
	Secret string `json:"secret"`
	Token  string `json:"token"`
}

type RequestDownloadResponse struct {
	Url      string `json:"url"`
	Password string `json:"password"`
	FileName string `json:"file_name"`
}

type DownloadResponse struct {
	Url string `json:"url"`
}
