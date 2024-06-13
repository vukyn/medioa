// Lưu ý ở tầng service:
// 	+ Hàm GetListPaging bổ sung config (0: get all, 1: chỉ get data, 2: chỉ get count).
// 	+ Tách thành hàm riêng những chức năng có khả năng tái sử dụng.
// Phải kèm theo lỗi error của hệ thống khi trả lỗi.

package service

import (
	"context"
	"medioa/config"
	"medioa/constants"
	"medioa/internal/storage/entity"
	"medioa/internal/storage/models"
	repo "medioa/internal/storage/repository"
	commonModel "medioa/models"
	"medioa/pkg/log"
	"medioa/pkg/routine"

)

type service struct {
	cfg    *config.Config
	lib    *commonModel.Lib
	repo   repo.IRepository
}

func InitService(cfg *config.Config, lib *commonModel.Lib, repo repo.IRepository) IService {
	return &service{
		cfg:  cfg,
		lib:  lib,
		repo: repo,
	}
}

func (s *service) GetById(ctx context.Context, id int64) (*models.Response, error) {
	log := log.New("service", "GetById")
	record, err := s.repo.GetById(ctx, id)
	if err != nil {
		log.Error("service.repo.GetById", err)
		return nil, err
	}
	if record == nil {
		return nil, nil
	}
	return record.Export(), nil
}

func (s *service) GetList(ctx context.Context, params *models.RequestParams) ([]*models.Response, error) {
	log := log.New("service", "GetList")
	queries := params.ToMap()
	records, err := s.repo.GetList(ctx, queries)
	if err != nil {
		log.Error("service.repo.GetList", err)
		return nil, err
	}
	return (&entity.Storage{}).ExportList(records), nil
}

func (s *service) GetListPaging(ctx context.Context, params *models.RequestParams) (*models.ListPaging, error) {
	log := log.New("service", "GetListPaging")
	queries := params.ToMap()
	errCh := make(chan error, 2)
	chCount := make(chan int64)
	chRecords := make(chan []*entity.Storage)
	if params.ConfigQuery == constants.CONFIG_QUERY_GET_ALL || params.ConfigQuery == constants.CONFIG_QUERY_GET_LIST {
		routine.Run(func() {
			records, err := s.repo.GetListPaging(ctx, queries)
			if err != nil {
				log.Error("service.repo.GetListPaging", err)
				errCh <- err
			}
			chRecords <- records
		})
	} else {
		routine.Run(func() {
			chRecords <- nil
		})
	}
	if params.ConfigQuery == constants.CONFIG_QUERY_GET_ALL || params.ConfigQuery == constants.CONFIG_QUERY_GET_COUNT {
		routine.Run(func() {
			count, err := s.repo.Count(ctx, queries)
			if err != nil {
				log.Error("service.repo.Count", err)
				errCh <- err
			}
			chCount <- count
		})
	} else {
		routine.Run(func() {
			chCount <- 0
		})
	}
	count := <-chCount
	records := <-chRecords
	close(errCh)
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return &models.ListPaging{
		ListPaging: commonModel.ListPaging{
			Page:  params.Page,
			Size:  params.Size,
			Count: count,
		},
		Records: (&entity.Storage{}).ExportList(records),
	}, nil
}

func (s *service) GetOne(ctx context.Context, params *models.RequestParams) (*models.Response, error) {
	log := log.New("service", "GetOne")
	queries := params.ToMap()
	record, err := s.repo.GetOne(ctx, queries)
	if err != nil {
		log.Error("service.repo.GetOne", err)
		return nil, err
	}
	if record == nil {
		return nil, nil
	}
	return record.Export(), nil
}

func (s *service) Count(ctx context.Context, params *models.RequestParams) (int64, error) {
	log := log.New("service", "Count")
	queries := params.ToMap()
	count, err := s.repo.Count(ctx, queries)
	if err != nil {
		log.Error("service.repo.Count", err)
		return 0, err
	}
	return count, nil
}

func (s *service) Create(ctx context.Context, userId int64, params *models.SaveRequest) (*models.Response, error) {
	log := log.New("service", "Create")
	obj := &entity.Storage{}
	obj.ParseForCreate(params, userId)
	res, err := s.repo.Create(ctx, obj)
	if err != nil {
		log.Error("service.repo.Create", err)
		return nil, err
	}
	return res.Export(), nil
}

func (s *service) CreateMany(ctx context.Context, userId int64, params []*models.SaveRequest) ([]*models.Response, error) {
	log := log.New("service", "CreateMany")
	objs := (&entity.Storage{}).ParseForCreateMany(params, userId)
	res, err := s.repo.CreateMany(ctx, objs)
	if err != nil {
		log.Error("service.repo.Create", err)
		return nil, err
	}
	return (&entity.Storage{}).ExportList(res), nil
}

func (s *service) Update(ctx context.Context, userId int64, params *models.SaveRequest) (*models.Response, error) {
	log := log.New("service", "Update")
	obj := &entity.Storage{}
	obj.ParseForUpdate(params, userId)
	res, err := s.repo.Update(ctx, obj)
	if err != nil {
		log.Error("service.repo.Update", err)
		return nil, err
	}
	return res.Export(), nil
}

func (s *service) UpdateMany(ctx context.Context, userId int64, params []*models.SaveRequest) (int64, error) {
	log := log.New("service", "UpdateMany")
	objs := (&entity.Storage{}).ParseForUpdateMany(params, userId)
	res, err := s.repo.UpdateMany(ctx, objs)
	if err != nil {
		log.Error("service.repo.UpdateMany", err)
		return 0, err
	}
	return res, nil
}

func (s *service) Upsert(ctx context.Context, userId int64, params *models.SaveRequest) (*models.Response, error) {
	res, err := s.GetById(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return s.Create(ctx, userId, params)
	}
	if _, err := s.Update(ctx, userId, params); err != nil {
		return nil, err
	}
	return s.GetById(ctx, params.Id)
}
