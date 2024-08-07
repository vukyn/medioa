package models

import (
	"medioa/pkg/xtype"
)

type UploadBlobRequest struct {
	SessionId string
	SecretId  string
	File      *xtype.File
}

type UploadBlobResponse struct {
	Url      string
	Token    string
	Ext      string
	FileName string
}

type UploadChunkRequest struct {
	SessionId   string
	Token       string
	FileName    string
	ChunkIndex  int64
	TotalChunks int64
	Chunk       *xtype.File
}

type UploadChunkResponse struct {
	Url      string
	Token    string
	BlockId  string
	Ext      string
	FileName string
}

type CommitChunkRequest struct {
	SessionId string
	Token     string
	FileName  string
	BlockIds  []string
}

type CommitChunkRsponse struct {
	TotalBlock int64
	FileSize   int64
}

type DownloadSASRequest struct {
	FileName string
}

type DownloadSASResponse struct {
	Url string
}
