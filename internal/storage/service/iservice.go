package service

import (
	"context"
	"medioa/internal/storage/models"
)

type IService interface {
	GetById(ctx context.Context, id int64) (*models.Response, error)
	GetOne(ctx context.Context, params *models.RequestParams) (*models.Response, error)
	GetList(ctx context.Context, params *models.RequestParams) ([]*models.Response, error)
	GetListPaging(ctx context.Context, params *models.RequestParams) (*models.ListPaging, error)
	Count(ctx context.Context, params *models.RequestParams) (int64, error)
	Create(ctx context.Context, userId int64, params *models.SaveRequest) (*models.Response, error)
	CreateMany(ctx context.Context, userId int64, params []*models.SaveRequest) ([]*models.Response, error)
	Update(ctx context.Context, userId int64, params *models.SaveRequest) (*models.Response, error)
	UpdateMany(ctx context.Context, userId int64, params []*models.SaveRequest) (int64, error)
	Upsert(ctx context.Context, userId int64, params *models.SaveRequest) (*models.Response, error)

	// Blob
	UploadPublicBlob(ctx context.Context, req *models.UploadBlobRequest) (*models.UploadBlobResponse, error)
	UploadPrivateBlob(ctx context.Context, req *models.UploadBlobRequest) (*models.UploadBlobResponse, error)
	DownloadSAS(ctx context.Context, req *models.DownloadSASRequest) (*models.DownloadSASResponse, error)
}
