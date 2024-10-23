package repository

import (
	"context"
	"medioa/internal/storage/entity"
)

type IRepository interface {
	GetById(ctx context.Context, id int64) (*entity.Storage, error)
	GetOne(ctx context.Context, queries map[string]any) (*entity.Storage, error)
	GetList(ctx context.Context, queries map[string]any) ([]*entity.Storage, error)
	GetListPaging(ctx context.Context, queries map[string]any) ([]*entity.Storage, error)
	Count(ctx context.Context, queries map[string]any) (int64, error)
	Create(ctx context.Context, obj *entity.Storage) (*entity.Storage, error)
	CreateMany(ctx context.Context, objs []*entity.Storage) ([]*entity.Storage, error)
	Update(ctx context.Context, obj *entity.Storage) (*entity.Storage, error)
	UpdateMany(ctx context.Context, objs []*entity.Storage) (int64, error)
}
