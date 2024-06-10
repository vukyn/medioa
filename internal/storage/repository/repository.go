// Kiểm tra validate cho hàm GET ONE, UPDATE ở tầng repository.
// Nên sử dụng subquery trong repository.
// Lưu ý ở tầng repository:
// 	+ Không 'join' vào những bảng khác một cách tùy tiện.
// 	+ Không tự ý dùng 'Group by' trong 'join'.
// 	+ Hạn chế sử dụng search query 'OR' để filter. Nên format lại data để search.

package repository

import (
	"context"
	"fmt"
	"medioa/constants"
	"medioa/internal/storage/entity"
	commonModel "medioa/models"

	"github.com/vukyn/kuery/conversion"
	"gorm.io/gorm"
)

type repo struct {
	lib       *commonModel.Lib
	tableName string
}

func InitRepo(lib *commonModel.Lib) IRepository {
	return &repo{
		lib:       lib,
		tableName: (&entity.Storage{}).TableName(),
	}
}

func (r *repo) dbWithContext(ctx context.Context) *gorm.DB {
	return r.lib.Db.WithContext(ctx)
}

func (r *repo) Create(ctx context.Context, obj *entity.Storage) (*entity.Storage, error) {
	result := r.dbWithContext(ctx).Create(obj)
	if result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *repo) CreateMany(ctx context.Context, objs []*entity.Storage) ([]*entity.Storage, error) {
	result := r.dbWithContext(ctx).Create(objs)
	if result.Error != nil {
		return nil, result.Error
	}
	return objs, nil
}

func (r *repo) Update(ctx context.Context, obj *entity.Storage) (*entity.Storage, error) {
	result := r.dbWithContext(ctx).Updates(obj)
	if result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *repo) UpdateMany(ctx context.Context, objs []*entity.Storage) (int64, error) {
	tx := r.dbWithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, obj := range objs {
		if err := tx.Updates(obj).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}
	return int64(len(objs)), nil
}

func (r *repo) initQuery(ctx context.Context, queries map[string]interface{}) *gorm.DB {
	obj := &entity.Storage{}
	query := r.dbWithContext(ctx).Model(obj)
	query = r.join(query, queries)
	query = r.column(query, queries)
	query = r.filter(query, queries)
	query = r.sort(query, queries)
	return query
}

func (r *repo) initSubQuery(ctx context.Context, queries map[string]interface{}) *gorm.DB {
	obj := &entity.Storage{}
	query := r.dbWithContext(ctx).Model(obj)
	query = r.joinSub(query, queries)
	query = r.columnSub(query, queries)
	query = r.filter(query, queries)
	query = r.sort(query, queries)
	return query
}

func (r *repo) getListGorm(ctx context.Context, queries map[string]interface{}) *gorm.DB {
	query := r.initQuery(ctx, queries)
	return query
}

func (r *repo) join(query *gorm.DB, queries map[string]interface{}) *gorm.DB {
	return query
}

func (r *repo) joinSub(query *gorm.DB, queries map[string]interface{}) *gorm.DB {
	return query
}

func (r *repo) column(query *gorm.DB, queries map[string]interface{}) *gorm.DB {
	query = query.Select(
		"storage.*",
	)
	return query
}

func (r *repo) columnSub(query *gorm.DB, queries map[string]interface{}) *gorm.DB {
	query = query.Select(
		"storage." + constants.FIELD_STORAGE_ID,
	)
	return query
}

func (r *repo) Count(ctx context.Context, queries map[string]interface{}) (int64, error) {
	var count int64
	if err := r.initSubQuery(ctx, queries).Select("count(1)").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repo) GetById(ctx context.Context, id int64) (*entity.Storage, error) {
	queries := map[string]interface{}{constants.FIELD_STORAGE_ID: id}
	return r.GetOne(ctx, queries)
}

func (r *repo) GetOne(ctx context.Context, queries map[string]interface{}) (*entity.Storage, error) {
	queries["size"] = int64(1)
	queries["page"] = int64(1)
	result, err := r.GetListPaging(ctx, queries)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result[0], nil
}

func (r *repo) GetList(ctx context.Context, queries map[string]interface{}) ([]*entity.Storage, error) {
	objs := []*entity.Storage{}
	query := r.getListGorm(ctx, queries)
	if err := query.Scan(&objs).Error; err != nil {
		return nil, err
	}
	return objs, nil
}

func (r *repo) GetListPaging(ctx context.Context, queries map[string]interface{}) ([]*entity.Storage, error) {
	objs := []*entity.Storage{}

	page := conversion.ReadInterfaceV2(queries, constants.FIELD_PAGE, constants.DEFAULT_PAGE, true)
	size := conversion.ReadInterfaceV2(queries, constants.FIELD_SIZE, constants.DEFAULT_SIZE, true)

	subQuery := r.initSubQuery(ctx, queries)
	subQuery = subQuery.Offset(int((page - 1) * size)).Limit(int(size))

	query := r.dbWithContext(ctx).Model(&entity.Storage{})
	query = r.join(query, queries)
	query = r.column(query, queries)
	query = r.sort(query, queries)
	query = query.Joins(
		fmt.Sprintf("inner join (?) as tmp on (%[1]s.%[2]s = tmp.%[2]s)", r.tableName, constants.FIELD_STORAGE_ID),
		subQuery,
	)

	if err := query.Scan(&objs).Error; err != nil {
		return nil, err
	}
	return objs, nil
}

func (r *repo) sort(query *gorm.DB, queries map[string]interface{}) *gorm.DB {
	sortMultiple := conversion.ReadInterfaceV2(queries, constants.FIELD_SORT_MULTIPLE, "", true)
	if sortMultiple != "" {
		query = query.Order(sortMultiple)
	} else {
		sortBy := conversion.ReadInterfaceV2(queries, constants.FIELD_SORT_BY, "", true)
		orderBy := conversion.ReadInterfaceV2(queries, constants.FIELD_ORDER_BY, constants.DEFAULT_SORT_ORDER, true)
		switch sortBy {
		case constants.FIELD_STORAGE_TYPE:
			query = query.Order(r.tableName + "." + constants.FIELD_STORAGE_TYPE + " " + orderBy)
		case constants.FIELD_STORAGE_TOKEN:
			query = query.Order(r.tableName + "." + constants.FIELD_STORAGE_TOKEN + " " + orderBy)
		case constants.FIELD_STORAGE_LIFE_TIME:
			query = query.Order(r.tableName + "." + constants.FIELD_STORAGE_LIFE_TIME + " " + orderBy)
		case constants.FIELD_STORAGE_CREATED_BY:
			query = query.Order(r.tableName + "." + constants.FIELD_STORAGE_CREATED_BY + " " + orderBy)
		case constants.FIELD_STORAGE_CREATED_AT:
			query = query.Order(r.tableName + "." + constants.FIELD_STORAGE_CREATED_AT + " " + orderBy)
		case constants.FIELD_STORAGE_ID:
			query = query.Order(r.tableName + "." + constants.FIELD_STORAGE_ID + " " + orderBy)
		case constants.FIELD_STORAGE_DOWNLOAD_URL:
			query = query.Order(r.tableName + "." + constants.FIELD_STORAGE_DOWNLOAD_URL + " " + orderBy)
		default:
			query = query.Order(r.tableName + "." + constants.FIELD_STORAGE_ID + " " + constants.SORT_ORDER_DESC)
		}
	}
	return query
}

func (r *repo) filter(query *gorm.DB, queries map[string]interface{}) *gorm.DB {
	id := conversion.ReadInterfaceV2(queries, constants.FIELD_STORAGE_ID, 0, true)
	downloadUrl := conversion.ReadInterfaceV2(queries, constants.FIELD_STORAGE_DOWNLOAD_URL, "", true)
	_type := conversion.ReadInterfaceV2(queries, constants.FIELD_STORAGE_TYPE, "", true)
	token := conversion.ReadInterfaceV2(queries, constants.FIELD_STORAGE_TOKEN, "", true)
	lifeTime := conversion.ReadInterfaceV2(queries, constants.FIELD_STORAGE_LIFE_TIME, 0, true)
	createdBy := conversion.ReadInterfaceV2(queries, constants.FIELD_STORAGE_CREATED_BY, 0, true)

	if id != 0 {
		query = query.Where(r.tableName+"."+constants.FIELD_STORAGE_ID+" = ? ", id)
	}
	if downloadUrl != "" {
		query = query.Where(r.tableName+"."+constants.FIELD_STORAGE_DOWNLOAD_URL+" = ? ", downloadUrl)
	}
	if _type != "" {
		query = query.Where(r.tableName+"."+constants.FIELD_STORAGE_TYPE+" = ? ", _type)
	}
	if token != "" {
		query = query.Where(r.tableName+"."+constants.FIELD_STORAGE_TOKEN+" = ? ", token)
	}
	if lifeTime != 0 {
		query = query.Where(r.tableName+"."+constants.FIELD_STORAGE_LIFE_TIME+" = ? ", lifeTime)
	}
	if createdBy != 0 {
		query = query.Where(r.tableName+"."+constants.FIELD_STORAGE_CREATED_BY+" = ? ", createdBy)
	}
	return query
}
