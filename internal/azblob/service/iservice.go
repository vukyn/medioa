package service

import (
	"context"
	"medioa/internal/azblob/models"
)

type IService interface {
	UploadPublicURL(ctx context.Context, req *models.UploadURLRequest) (*models.UploadResponse, error)
	UploadPublicBlob(ctx context.Context, req *models.UploadBlobRequest) (*models.UploadResponse, error)
	UploadPublicChunk(ctx context.Context, req *models.UploadChunkRequest) (*models.UploadChunkResponse, error)
	CommitPublicChunk(ctx context.Context, req *models.CommitChunkRequest) (*models.CommitChunkRsponse, error)
	UploadPrivateBlob(ctx context.Context, req *models.UploadBlobRequest) (*models.UploadResponse, error)
	UploadPrivateChunk(ctx context.Context, req *models.UploadChunkRequest) (*models.UploadChunkResponse, error)
	CommitPrivateChunk(ctx context.Context, req *models.CommitChunkRequest) (*models.CommitChunkRsponse, error)
	DownloadSAS(ctx context.Context, req *models.DownloadSASRequest) (*models.DownloadSASResponse, error)
}
