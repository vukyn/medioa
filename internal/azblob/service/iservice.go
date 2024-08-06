package service

import (
	"context"
	"medioa/internal/azblob/models"
)

type IService interface {
	UploadPublicBlob(ctx context.Context, req *models.UploadBlobRequest) (*models.UploadBlobResponse, error)
	UploadPublicChunk(ctx context.Context, req *models.UploadChunkRequest) (*models.UploadChunkResponse, error)
	CommitPublicChunk(ctx context.Context, req *models.CommitChunkRequest) error
	UploadPrivateBlob(ctx context.Context, req *models.UploadBlobRequest) (*models.UploadBlobResponse, error)
	DownloadSAS(ctx context.Context, req *models.DownloadSASRequest) (*models.DownloadSASResponse, error)
}
