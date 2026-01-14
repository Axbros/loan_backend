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

var _ LoanUserSmsRecordsDao = (*loanUserSmsRecordsDao)(nil)

// LoanUserSmsRecordsDao defining the dao interface
type LoanUserSmsRecordsDao interface {
	Create(ctx context.Context, table *model.LoanUserSmsRecords) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.LoanUserSmsRecords) error
	GetByID(ctx context.Context, id uint64) (*model.LoanUserSmsRecords, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanUserSmsRecords, int64, error)

	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.LoanUserSmsRecords, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanUserSmsRecords, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanUserSmsRecords, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanUserSmsRecords) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanUserSmsRecords) error
}

type loanUserSmsRecordsDao struct {
	db    *gorm.DB
	cache cache.LoanUserSmsRecordsCache // if nil, the cache is not used.
	sfg   *singleflight.Group           // if cache is nil, the sfg is not used.
}

// NewLoanUserSmsRecordsDao creating the dao interface
func NewLoanUserSmsRecordsDao(db *gorm.DB, xCache cache.LoanUserSmsRecordsCache) LoanUserSmsRecordsDao {
	if xCache == nil {
		return &loanUserSmsRecordsDao{db: db}
	}
	return &loanUserSmsRecordsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *loanUserSmsRecordsDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new loanUserSmsRecords, insert the record and the id value is written back to the table
func (d *loanUserSmsRecordsDao) Create(ctx context.Context, table *model.LoanUserSmsRecords) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a loanUserSmsRecords by id
func (d *loanUserSmsRecordsDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.LoanUserSmsRecords{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a loanUserSmsRecords by ids
func (d *loanUserSmsRecordsDao) UpdateByID(ctx context.Context, table *model.LoanUserSmsRecords) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *loanUserSmsRecordsDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.LoanUserSmsRecords) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.BaseinfoID != 0 {
		update["baseinfo_id"] = table.BaseinfoID
	}
	if table.Direction != 0 {
		update["direction"] = table.Direction
	}
	if table.Address != "" {
		update["address"] = table.Address
	}
	if table.SmsTime != nil && table.SmsTime.IsZero() == false {
		update["sms_time"] = table.SmsTime
	}
	if table.Body != "" {
		update["body"] = table.Body
	}
	if table.BodyHash != "" {
		update["body_hash"] = table.BodyHash
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a loanUserSmsRecords by id
func (d *loanUserSmsRecordsDao) GetByID(ctx context.Context, id uint64) (*model.LoanUserSmsRecords, error) {
	// no cache
	if d.cache == nil {
		record := &model.LoanUserSmsRecords{}
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
			table := &model.LoanUserSmsRecords{}
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
			if err = d.cache.Set(ctx, id, table, cache.LoanUserSmsRecordsExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.LoanUserSmsRecords)
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

// GetByColumns get a paginated list of loanUserSmsRecordss by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *loanUserSmsRecordsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanUserSmsRecords, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.LoanUserSmsRecordsColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.LoanUserSmsRecords{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.LoanUserSmsRecords{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// DeleteByIDs batch delete loanUserSmsRecords by ids
func (d *loanUserSmsRecordsDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.LoanUserSmsRecords{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.deleteCache(ctx, id)
	}

	return nil
}

// GetByCondition get a loanUserSmsRecords by custom condition
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html#_2-condition-parameters-optional
func (d *loanUserSmsRecordsDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.LoanUserSmsRecords, error) {
	queryStr, args, err := c.ConvertToGorm(query.WithWhitelistNames(model.LoanUserSmsRecordsColumnNames))
	if err != nil {
		return nil, err
	}

	table := &model.LoanUserSmsRecords{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs Batch get loanUserSmsRecords by ids
func (d *loanUserSmsRecordsDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanUserSmsRecords, error) {
	// no cache
	if d.cache == nil {
		var records []*model.LoanUserSmsRecords
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.LoanUserSmsRecords)
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
			var records []*model.LoanUserSmsRecords
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
				if err = d.cache.MultiSet(ctx, records, cache.LoanUserSmsRecordsExpireTime); err != nil {
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

// GetByLastID Get a paginated list of loanUserSmsRecordss by last id
func (d *loanUserSmsRecordsDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanUserSmsRecords, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.LoanUserSmsRecords{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *loanUserSmsRecordsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanUserSmsRecords) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *loanUserSmsRecordsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.LoanUserSmsRecords{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *loanUserSmsRecordsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanUserSmsRecords) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
