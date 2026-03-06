package dao

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"loan/internal/cache"
	"loan/internal/database"
	"loan/internal/model"
)

var _ LoanCollectionLogsDao = (*loanCollectionLogsDao)(nil)

// LoanCollectionLogsDao defining the dao interface
type LoanCollectionLogsDao interface {
	Create(ctx context.Context, table *model.LoanCollectionLogs) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.LoanCollectionLogs) error
	GetByID(ctx context.Context, id uint64) (*model.LoanCollectionLogs, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanCollectionLogs, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanCollectionLogs) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanCollectionLogs) error
}

type loanCollectionLogsDao struct {
	db    *gorm.DB
	cache cache.LoanCollectionLogsCache // if nil, the cache is not used.
	sfg   *singleflight.Group           // if cache is nil, the sfg is not used.
}

// NewLoanCollectionLogsDao creating the dao interface
func NewLoanCollectionLogsDao(db *gorm.DB, xCache cache.LoanCollectionLogsCache) LoanCollectionLogsDao {
	if xCache == nil {
		return &loanCollectionLogsDao{db: db}
	}
	return &loanCollectionLogsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *loanCollectionLogsDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new loanCollectionLogs, insert the record and the id value is written back to the table
func (d *loanCollectionLogsDao) Create(ctx context.Context, table *model.LoanCollectionLogs) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a loanCollectionLogs by id
func (d *loanCollectionLogsDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.LoanCollectionLogs{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a loanCollectionLogs by ids
func (d *loanCollectionLogsDao) UpdateByID(ctx context.Context, table *model.LoanCollectionLogs) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *loanCollectionLogsDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.LoanCollectionLogs) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.CaseID != 0 {
		update["case_id"] = table.CaseID
	}
	if table.CollectorUserID != 0 {
		update["collector_user_id"] = table.CollectorUserID
	}
	if table.ActionType != "" {
		update["action_type"] = table.ActionType
	}
	if table.Content != "" {
		update["content"] = table.Content
	}
	if table.NextFollowUpAt != nil && table.NextFollowUpAt.IsZero() == false {
		update["next_follow_up_at"] = table.NextFollowUpAt
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a loanCollectionLogs by id
func (d *loanCollectionLogsDao) GetByID(ctx context.Context, id uint64) (*model.LoanCollectionLogs, error) {
	// no cache
	if d.cache == nil {
		record := &model.LoanCollectionLogs{}
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
		val, err, _ := d.sfg.Do(utils.Uint64ToStr(id), func() (interface{}, error) {
			table := &model.LoanCollectionLogs{}
			err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
			if err != nil {
				// set placeholder cache to prevent cache penetration, default expiration time 10 minutes
				if errors.Is(err, database.ErrRecordNotFound) {
					if err = d.cache.SetPlaceholder(ctx, id); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("id", id))
					}
					return nil, database.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			if err = d.cache.Set(ctx, id, table, cache.LoanCollectionLogsExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.LoanCollectionLogs)
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

// GetByColumns get a paginated list of loanCollectionLogss by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *loanCollectionLogsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanCollectionLogs, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.LoanCollectionLogsColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.LoanCollectionLogs{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.LoanCollectionLogs{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *loanCollectionLogsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanCollectionLogs) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *loanCollectionLogsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.LoanCollectionLogs{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *loanCollectionLogsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanCollectionLogs) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
