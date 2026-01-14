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

var _ LoanLoginAuditDao = (*loanLoginAuditDao)(nil)

// LoanLoginAuditDao defining the dao interface
type LoanLoginAuditDao interface {
	Create(ctx context.Context, table *model.LoanLoginAudit) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.LoanLoginAudit) error
	GetByID(ctx context.Context, id uint64) (*model.LoanLoginAudit, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanLoginAudit, int64, error)

	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.LoanLoginAudit, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanLoginAudit, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanLoginAudit, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanLoginAudit) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanLoginAudit) error
}

type loanLoginAuditDao struct {
	db    *gorm.DB
	cache cache.LoanLoginAuditCache // if nil, the cache is not used.
	sfg   *singleflight.Group       // if cache is nil, the sfg is not used.
}

// NewLoanLoginAuditDao creating the dao interface
func NewLoanLoginAuditDao(db *gorm.DB, xCache cache.LoanLoginAuditCache) LoanLoginAuditDao {
	if xCache == nil {
		return &loanLoginAuditDao{db: db}
	}
	return &loanLoginAuditDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *loanLoginAuditDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new loanLoginAudit, insert the record and the id value is written back to the table
func (d *loanLoginAuditDao) Create(ctx context.Context, table *model.LoanLoginAudit) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a loanLoginAudit by id
func (d *loanLoginAuditDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.LoanLoginAudit{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a loanLoginAudit by ids
func (d *loanLoginAuditDao) UpdateByID(ctx context.Context, table *model.LoanLoginAudit) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *loanLoginAuditDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.LoanLoginAudit) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.UserID != 0 {
		update["user_id"] = table.UserID
	}
	if table.LoginType != "" {
		update["login_type"] = table.LoginType
	}
	if table.IP != "" {
		update["ip"] = table.IP
	}
	if table.UserAgent != "" {
		update["user_agent"] = table.UserAgent
	}
	if table.Success != 0 {
		update["success"] = table.Success
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a loanLoginAudit by id
func (d *loanLoginAuditDao) GetByID(ctx context.Context, id uint64) (*model.LoanLoginAudit, error) {
	// no cache
	if d.cache == nil {
		record := &model.LoanLoginAudit{}
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
			table := &model.LoanLoginAudit{}
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
			if err = d.cache.Set(ctx, id, table, cache.LoanLoginAuditExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.LoanLoginAudit)
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

// GetByColumns get a paginated list of loanLoginAudits by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *loanLoginAuditDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanLoginAudit, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.LoanLoginAuditColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.LoanLoginAudit{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.LoanLoginAudit{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// DeleteByIDs batch delete loanLoginAudit by ids
func (d *loanLoginAuditDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.LoanLoginAudit{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.deleteCache(ctx, id)
	}

	return nil
}

// GetByCondition get a loanLoginAudit by custom condition
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html#_2-condition-parameters-optional
func (d *loanLoginAuditDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.LoanLoginAudit, error) {
	queryStr, args, err := c.ConvertToGorm(query.WithWhitelistNames(model.LoanLoginAuditColumnNames))
	if err != nil {
		return nil, err
	}

	table := &model.LoanLoginAudit{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs Batch get loanLoginAudit by ids
func (d *loanLoginAuditDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanLoginAudit, error) {
	// no cache
	if d.cache == nil {
		var records []*model.LoanLoginAudit
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.LoanLoginAudit)
		for _, record := range records {
			itemMap[record.ID] = record
		}
		return itemMap, nil
	}

	// get form cache
	itemMap, err := d.cache.MultiGet(ctx, ids)
	if err != nil {
		return nil, err
	}

	var missedIDs []uint64
	for _, id := range ids {
		if _, ok := itemMap[id]; !ok {
			missedIDs = append(missedIDs, id)
		}
	}

	// get missed data
	if len(missedIDs) > 0 {
		// find the id of an active placeholder, i.e. an id that does not exist in database
		var realMissedIDs []uint64
		for _, id := range missedIDs {
			_, err = d.cache.Get(ctx, id)
			if d.cache.IsPlaceholderErr(err) {
				continue
			}
			realMissedIDs = append(realMissedIDs, id)
		}

		// get missed id from database
		if len(realMissedIDs) > 0 {
			var records []*model.LoanLoginAudit
			var recordIDMap = make(map[uint64]struct{})
			err = d.db.WithContext(ctx).Where("id IN (?)", realMissedIDs).Find(&records).Error
			if err != nil {
				return nil, err
			}
			if len(records) > 0 {
				for _, record := range records {
					itemMap[record.ID] = record
					recordIDMap[record.ID] = struct{}{}
				}
				if err = d.cache.MultiSet(ctx, records, cache.LoanLoginAuditExpireTime); err != nil {
					logger.Warn("cache.MultiSet error", logger.Err(err), logger.Any("ids", records))
				}
				if len(records) == len(realMissedIDs) {
					return itemMap, nil
				}
			}
			for _, id := range realMissedIDs {
				if _, ok := recordIDMap[id]; !ok {
					if err = d.cache.SetPlaceholder(ctx, id); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("id", id))
					}
				}
			}
		}
	}

	return itemMap, nil
}

// GetByLastID Get a paginated list of loanLoginAudits by last id
func (d *loanLoginAuditDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanLoginAudit, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.LoanLoginAudit{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *loanLoginAuditDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanLoginAudit) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *loanLoginAuditDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.LoanLoginAudit{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *loanLoginAuditDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanLoginAudit) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
