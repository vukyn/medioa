package repository

import (
	"context"
	"medioa/internal/secret/entity"
)

type IRepository interface {
	GetById(ctx context.Context, id int64) (*entity.Secret, error)
	GetOne(ctx context.Context, queries map[string]any) (*entity.Secret, error)
	GetList(ctx context.Context, queries map[string]any) ([]*entity.Secret, error)
	GetListPaging(ctx context.Context, queries map[string]any) ([]*entity.Secret, error)
	Count(ctx context.Context, queries map[string]any) (int64, error)
	Create(ctx context.Context, obj *entity.Secret) (*entity.Secret, error)
	CreateMany(ctx context.Context, objs []*entity.Secret) ([]*entity.Secret, error)
	Update(ctx context.Context, obj *entity.Secret) (*entity.Secret, error)
	UpdateMany(ctx context.Context, objs []*entity.Secret) (int64, error)
}
