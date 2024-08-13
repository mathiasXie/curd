package curd

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Model[T schema.Tabler] struct {
	Db       []*gorm.DB
	isMaster bool
}

type ModelBO struct {
	PageNo     int32
	PageSize   int32
	Limit      int32
	Offset     int32
	Orders     string
	Conditions []map[string]interface{}
}

// NewModel 可以给出多个DB实例，默认第一个作为主库也就是写库，其他作为读库
func NewModel[T schema.Tabler](db ...*gorm.DB) *Model[T] {
	return &Model[T]{
		Db: db,
	}
}

func (model *Model[T]) IsMaster(isMaster bool) *Model[T] {
	model.isMaster = isMaster
	return model
}

func (model *Model[T]) GetDb(ctx context.Context) *gorm.DB {
	if len(model.Db) == 1 || model.isMaster {
		return model.Db[0].WithContext(ctx)
	}
	return model.Db[1].WithContext(ctx)
}

func (model *Model[T]) SelectAll(ctx context.Context) ([]*T, int64, error) {
	var result []*T
	err := model.GetDb(ctx).Model(new(T)).Find(&result).Error
	if err != nil {
		return nil, 0, err
	}
	return result, int64(len(result)), nil
}

func (model *Model[T]) Select(ctx context.Context, bo *ModelBO) ([]*T, int64, error) {
	var result []*T
	db := model.GetDb(ctx).Model(new(T))
	db = model.addQueryConditions(db, bo)
	db = model.addLimit(db, bo)
	db = model.addOrder(db, bo)
	err := db.Find(&result).Error
	if err != nil {
		return nil, 0, err
	}
	return result, int64(len(result)), nil
}

func (model *Model[T]) FindOne(ctx context.Context, id int64) (*T, error) {
	var result T
	err := model.GetDb(ctx).Where("id = ?", id).First(&result).Error
	return &result, err
}

func (model *Model[T]) Insert(ctx context.Context, data map[string]interface{}) (*T, error) {
	err := model.GetDb(ctx).Model(new(T)).Create(&data).Error
	if err != nil {
		return nil, err
	}
	var m *T
	model.GetDb(ctx).Last(&m)
	return m, nil
}

func (model *Model[T]) Update(ctx context.Context, id int64, data map[string]interface{}) error {
	return model.GetDb(ctx).Model(new(T)).Where("id=?", id).Updates(data).Error
}

func (model *Model[T]) Delete(ctx context.Context, data *T) error {
	return model.GetDb(ctx).Delete(data).Error
}

// 查询条件
// 筛选条件(列表元素间OR,Map元素间AND
func (model *Model[T]) addQueryConditions(db *gorm.DB, bo *ModelBO) *gorm.DB {
	// 查询条件
	orFlag := false // 是否用OR运算符
	for _, condition := range bo.Conditions {

		for k, v := range condition {
			if orFlag {
				db = db.Or(k, v) // ignore_security_alert
				orFlag = false
			}
			db = db.Where(k, v) // ignore_security_alert
		}
		orFlag = true
	}
	return db
}

// 分页+条数限制
func (model *Model[T]) addLimit(db *gorm.DB, bo *ModelBO) *gorm.DB {
	// 分页查询
	if bo.PageNo > 0 && bo.PageSize > 0 {
		db = db.Offset(int((bo.PageNo - 1) * bo.PageSize))
		db = db.Limit(int(bo.PageSize))
	}
	// 条数限制
	//if limit > 0 {
	//	db = db.Limit(int(limit))
	//}
	return db
}

// 排序
func (model *Model[T]) addOrder(db *gorm.DB, bo *ModelBO) *gorm.DB {
	if db == nil || bo == nil {
		return db
	}
	// 排序
	for _, order := range bo.Orders {
		db = db.Order(order)
	}
	return db
}
