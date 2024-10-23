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
	"medioa/internal/secret/entity"
	commonModel "medioa/models"

	"github.com/vukyn/kuery/conv"
	"gorm.io/gorm"
)

type repo struct {
	lib       *commonModel.Lib
	tableName string
}

func InitRepo(lib *commonModel.Lib) IRepository {
	return &repo{
		lib:       lib,
		tableName: (&entity.Secret{}).TableName(),
	}
}

func (r *repo) dbWithContext(ctx context.Context) *gorm.DB {
	return r.lib.Db.WithContext(ctx)
}

func (r *repo) Create(ctx context.Context, obj *entity.Secret) (*entity.Secret, error) {
	result := r.dbWithContext(ctx).Create(obj)
	if result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *repo) CreateMany(ctx context.Context, objs []*entity.Secret) ([]*entity.Secret, error) {
	result := r.dbWithContext(ctx).Create(objs)
	if result.Error != nil {
		return nil, result.Error
	}
	return objs, nil
}

func (r *repo) Update(ctx context.Context, obj *entity.Secret) (*entity.Secret, error) {
	result := r.dbWithContext(ctx).Updates(obj)
	if result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *repo) UpdateMany(ctx context.Context, objs []*entity.Secret) (int64, error) {
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

func (r *repo) initQuery(ctx context.Context, queries map[string]any) *gorm.DB {
	obj := &entity.Secret{}
	query := r.dbWithContext(ctx).Model(obj)
	query = r.join(query, queries)
	query = r.column(query, queries)
	query = r.filter(query, queries)
	query = r.sort(query, queries)
	return query
}

func (r *repo) initSubQuery(ctx context.Context, queries map[string]any) *gorm.DB {
	obj := &entity.Secret{}
	query := r.dbWithContext(ctx).Model(obj)
	query = r.joinSub(query, queries)
	query = r.columnSub(query, queries)
	query = r.filter(query, queries)
	query = r.sort(query, queries)
	return query
}

func (r *repo) getListGorm(ctx context.Context, queries map[string]any) *gorm.DB {
	query := r.initQuery(ctx, queries)
	return query
}

func (r *repo) join(query *gorm.DB, queries map[string]any) *gorm.DB {
	return query
}

func (r *repo) joinSub(query *gorm.DB, queries map[string]any) *gorm.DB {
	return query
}

func (r *repo) column(query *gorm.DB, queries map[string]any) *gorm.DB {
	query = query.Select(
		"secrets.*",
	)
	return query
}

func (r *repo) columnSub(query *gorm.DB, queries map[string]any) *gorm.DB {
	query = query.Select(
		"secrets." + constants.FIELD_SECRET_ID,
	)
	return query
}

func (r *repo) Count(ctx context.Context, queries map[string]any) (int64, error) {
	var count int64
	if err := r.initSubQuery(ctx, queries).Select("count(1)").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repo) GetById(ctx context.Context, id int64) (*entity.Secret, error) {
	queries := map[string]any{constants.FIELD_SECRET_ID: id}
	return r.GetOne(ctx, queries)
}

func (r *repo) GetOne(ctx context.Context, queries map[string]any) (*entity.Secret, error) {
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

func (r *repo) GetList(ctx context.Context, queries map[string]any) ([]*entity.Secret, error) {
	objs := []*entity.Secret{}
	query := r.getListGorm(ctx, queries)
	if err := query.Scan(&objs).Error; err != nil {
		return nil, err
	}
	return objs, nil
}

func (r *repo) GetListPaging(ctx context.Context, queries map[string]any) ([]*entity.Secret, error) {
	objs := []*entity.Secret{}

	page := conv.ReadInterface(queries, constants.FIELD_PAGE, constants.DEFAULT_PAGE)
	size := conv.ReadInterface(queries, constants.FIELD_SIZE, constants.DEFAULT_SIZE)

	subQuery := r.initSubQuery(ctx, queries)
	subQuery = subQuery.Offset(int((page - 1) * size)).Limit(int(size))

	query := r.dbWithContext(ctx).Model(&entity.Secret{})
	query = r.join(query, queries)
	query = r.column(query, queries)
	query = r.sort(query, queries)
	query = query.Joins(
		fmt.Sprintf("inner join (?) as tmp on (%[1]s.%[2]s = tmp.%[2]s)", r.tableName, constants.FIELD_SECRET_ID),
		subQuery,
	)

	if err := query.Scan(&objs).Error; err != nil {
		return nil, err
	}
	return objs, nil
}

func (r *repo) sort(query *gorm.DB, queries map[string]any) *gorm.DB {
	sortMultiple := conv.ReadInterface(queries, constants.FIELD_SORT_MULTIPLE, "")
	if sortMultiple != "" {
		query = query.Order(sortMultiple)
	} else {
		sortBy := conv.ReadInterface(queries, constants.FIELD_SORT_BY, "")
		orderBy := conv.ReadInterface(queries, constants.FIELD_ORDER_BY, constants.DEFAULT_SORT_ORDER)
		switch sortBy {
		case constants.FIELD_SECRET_USERNAME:
			query = query.Order(r.tableName + "." + constants.FIELD_SECRET_USERNAME + " " + orderBy)
		case constants.FIELD_SECRET_PASSWORD:
			query = query.Order(r.tableName + "." + constants.FIELD_SECRET_PASSWORD + " " + orderBy)
		case constants.FIELD_SECRET_PIN_CODE:
			query = query.Order(r.tableName + "." + constants.FIELD_SECRET_PIN_CODE + " " + orderBy)
		case constants.FIELD_SECRET_CREATED_BY:
			query = query.Order(r.tableName + "." + constants.FIELD_SECRET_CREATED_BY + " " + orderBy)
		case constants.FIELD_SECRET_CREATED_AT:
			query = query.Order(r.tableName + "." + constants.FIELD_SECRET_CREATED_AT + " " + orderBy)
		case constants.FIELD_SECRET_ID:
			query = query.Order(r.tableName + "." + constants.FIELD_SECRET_ID + " " + orderBy)
		case constants.FIELD_SECRET_UUID:
			query = query.Order(r.tableName + "." + constants.FIELD_SECRET_UUID + " " + orderBy)
		default:
			query = query.Order(r.tableName + "." + constants.FIELD_SECRET_ID + " " + constants.SORT_ORDER_DESC)
		}
	}
	return query
}

func (r *repo) filter(query *gorm.DB, queries map[string]any) *gorm.DB {
	id := conv.ReadInterface(queries, constants.FIELD_SECRET_ID, int64(0))
	uuid := conv.ReadInterface(queries, constants.FIELD_SECRET_UUID, "")
	username := conv.ReadInterface(queries, constants.FIELD_SECRET_USERNAME, "")
	password := conv.ReadInterface(queries, constants.FIELD_SECRET_PASSWORD, "")
	pinCode := conv.ReadInterface(queries, constants.FIELD_SECRET_PIN_CODE, "")
	createdBy := conv.ReadInterface(queries, constants.FIELD_SECRET_CREATED_BY, int64(0))

	if id != int64(0) {
		query = query.Where(r.tableName+"."+constants.FIELD_SECRET_ID+" = ? ", id)
	}
	if uuid != "" {
		query = query.Where(r.tableName+"."+constants.FIELD_SECRET_UUID+" = ? ", uuid)
	}
	if username != "" {
		query = query.Where(r.tableName+"."+constants.FIELD_SECRET_USERNAME+" = ? ", username)
	}
	if password != "" {
		query = query.Where(r.tableName+"."+constants.FIELD_SECRET_PASSWORD+" = ? ", password)
	}
	if pinCode != "" {
		query = query.Where(r.tableName+"."+constants.FIELD_SECRET_PIN_CODE+" = ? ", pinCode)
	}
	if createdBy != int64(0) {
		query = query.Where(r.tableName+"."+constants.FIELD_SECRET_CREATED_BY+" = ? ", createdBy)
	}
	return query
}
