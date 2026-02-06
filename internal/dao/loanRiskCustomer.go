package dao

import (
	"context"
	"errors"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"loan/internal/cache"
	"loan/internal/database"
	"loan/internal/model"
)

var _ LoanRiskCustomerDao = (*loanRiskCustomerDao)(nil)

// LoanRiskCustomerDao defining the dao interface
type LoanRiskCustomerDao interface {
	Create(ctx context.Context, table *model.LoanRiskCustomer) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.LoanRiskCustomer) error
	GetByID(ctx context.Context, id uint64) (*model.LoanRiskCustomer, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanRiskCustomer, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanRiskCustomer) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanRiskCustomer) error
}

type loanRiskCustomerDao struct {
	db    *gorm.DB
	cache cache.LoanRiskCustomerCache // if nil, the cache is not used.
	sfg   *singleflight.Group         // if cache is nil, the sfg is not used.
}

// NewLoanRiskCustomerDao creating the dao interface
func NewLoanRiskCustomerDao(db *gorm.DB, xCache cache.LoanRiskCustomerCache) LoanRiskCustomerDao {
	if xCache == nil {
		return &loanRiskCustomerDao{db: db}
	}
	return &loanRiskCustomerDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *loanRiskCustomerDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new loanRiskCustomer, insert the record and the id value is written back to the table
func (d *loanRiskCustomerDao) Create(ctx context.Context, table *model.LoanRiskCustomer) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a loanRiskCustomer by id
func (d *loanRiskCustomerDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.LoanRiskCustomer{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a loanRiskCustomer by id, support partial update
func (d *loanRiskCustomerDao) UpdateByID(ctx context.Context, table *model.LoanRiskCustomer) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *loanRiskCustomerDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.LoanRiskCustomer) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.LoanBaseinfoID != 0 {
		update["loan_baseinfo_id"] = table.LoanBaseinfoID
	}
	if table.RiskType != 0 {
		update["risk_type"] = table.RiskType
	}
	if table.RiskReason != "" {
		update["risk_reason"] = table.RiskReason
	}
	if table.CreatedBy != 0 {
		update["created_by"] = table.CreatedBy
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a loanRiskCustomer by id
func (d *loanRiskCustomerDao) GetByID(ctx context.Context, id uint64) (*model.LoanRiskCustomer, error) {
	// no cache
	if d.cache == nil {
		record := &model.LoanRiskCustomer{}
		err := d.db.WithContext(ctx).Where("id = ?", id).First(record).Error
		return record, err
	}

	// get from cache
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	// get from database
	if errors.Is(err, database.ErrCacheNotFound) {
		// for the same id, prevent high concurrent simultaneous access to database
		val, err, _ := d.sfg.Do(utils.Uint64ToStr(id), func() (interface{}, error) { //nolint
			table := &model.LoanRiskCustomer{}
			err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
			if err != nil {
				if errors.Is(err, database.ErrRecordNotFound) {
					// set placeholder cache to prevent cache penetration, default expiration time 10 minutes
					if err = d.cache.SetPlaceholder(ctx, id); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("id", id))
					}
					return nil, database.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			if err = d.cache.Set(ctx, id, table, cache.LoanRiskCustomerExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.LoanRiskCustomer)
		if !ok {
			return nil, database.ErrRecordNotFound
		}
		return table, nil
	}

	if d.cache.IsPlaceholderErr(err) {
		return nil, database.ErrRecordNotFound
	}

	return nil, err
}

// GetByColumns get a paginated list of loanRiskCustomers by custom conditions.
// For more details, please refer to https://go-sponge.com/component/custom-page-query.html
func (d *loanRiskCustomerDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanRiskCustomer, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.LoanRiskCustomerColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	// 核心修改：同时Preload LoanBaseinfo + OperateUser（仅查id/username，性能优化）
	dbQuery := d.db.WithContext(ctx).Model(&model.LoanRiskCustomer{}).
		Preload("LoanBaseinfo"). // 预加载贷款基础信息完整数据
		Preload("OperateUser", func(db *gorm.DB) *gorm.DB {
			// 仅查询操作人核心字段，避免返回密码/部门等冗余敏感数据
			return db.Select("id, username")
		}).
		Where(queryStr, args...)

	var total int64
	if params.Sort != "ignore count" { // 总数统计无需预加载，单独查询提升性能
		err = d.db.WithContext(ctx).Model(&model.LoanRiskCustomer{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.LoanRiskCustomer{}
	order, limit, offset := params.ConvertToPage()
	// 复用预加载的查询器执行分页查询
	err = dbQuery.Order(order).Limit(limit).Offset(offset).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *loanRiskCustomerDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanRiskCustomer) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *loanRiskCustomerDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.LoanRiskCustomer{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *loanRiskCustomerDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanRiskCustomer) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
