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

var _ LoanPaymentChannelsDao = (*loanPaymentChannelsDao)(nil)

// LoanPaymentChannelsDao defining the dao interface
type LoanPaymentChannelsDao interface {
	Create(ctx context.Context, table *model.LoanPaymentChannels) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.LoanPaymentChannels) error
	GetByID(ctx context.Context, id uint64) (*model.LoanPaymentChannels, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanPaymentChannels, int64, error)

	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.LoanPaymentChannels, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanPaymentChannels, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanPaymentChannels, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanPaymentChannels) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanPaymentChannels) error
}

type loanPaymentChannelsDao struct {
	db    *gorm.DB
	cache cache.LoanPaymentChannelsCache // if nil, the cache is not used.
	sfg   *singleflight.Group            // if cache is nil, the sfg is not used.
}

// NewLoanPaymentChannelsDao creating the dao interface
func NewLoanPaymentChannelsDao(db *gorm.DB, xCache cache.LoanPaymentChannelsCache) LoanPaymentChannelsDao {
	if xCache == nil {
		return &loanPaymentChannelsDao{db: db}
	}
	return &loanPaymentChannelsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *loanPaymentChannelsDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new loanPaymentChannels, insert the record and the id value is written back to the table
func (d *loanPaymentChannelsDao) Create(ctx context.Context, table *model.LoanPaymentChannels) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a loanPaymentChannels by id
func (d *loanPaymentChannelsDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.LoanPaymentChannels{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a loanPaymentChannels by ids
func (d *loanPaymentChannelsDao) UpdateByID(ctx context.Context, table *model.LoanPaymentChannels) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *loanPaymentChannelsDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.LoanPaymentChannels) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.Code != "" {
		update["code"] = table.Code
	}
	if table.Name != "" {
		update["name"] = table.Name
	}
	if table.MerchantNo != "" {
		update["merchant_no"] = table.MerchantNo
	}
	if table.Status != 0 {
		update["status"] = table.Status
	}
	if table.CanPayout != 0 {
		update["can_payout"] = table.CanPayout
	}
	if table.CanCollect != 0 {
		update["can_collect"] = table.CanCollect
	}
	if table.PayoutFeeRate == 0.0 {
		update["payout_fee_rate"] = table.PayoutFeeRate
	}
	if table.PayoutFeeFixed != 0 {
		update["payout_fee_fixed"] = table.PayoutFeeFixed
	}
	if table.CollectFeeRate == 0.0 {
		update["collect_fee_rate"] = table.CollectFeeRate
	}
	if table.CollectFeeFixed != 0 {
		update["collect_fee_fixed"] = table.CollectFeeFixed
	}
	if table.CollectMinAmount != 0 {
		update["collect_min_amount"] = table.CollectMinAmount
	}
	if table.CollectMaxAmount != 0 {
		update["collect_max_amount"] = table.CollectMaxAmount
	}
	if table.PayoutMinAmount != 0 {
		update["payout_min_amount"] = table.PayoutMinAmount
	}
	if table.PayoutMaxAmount != 0 {
		update["payout_max_amount"] = table.PayoutMaxAmount
	}
	if table.SettlementCycle != "" {
		update["settlement_cycle"] = table.SettlementCycle
	}
	if table.SettlementDesc != "" {
		update["settlement_desc"] = table.SettlementDesc
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a loanPaymentChannels by id
func (d *loanPaymentChannelsDao) GetByID(ctx context.Context, id uint64) (*model.LoanPaymentChannels, error) {
	// no cache
	if d.cache == nil {
		record := &model.LoanPaymentChannels{}
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
			table := &model.LoanPaymentChannels{}
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
			if err = d.cache.Set(ctx, id, table, cache.LoanPaymentChannelsExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.LoanPaymentChannels)
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

// GetByColumns get a paginated list of loanPaymentChannelss by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *loanPaymentChannelsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanPaymentChannels, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.LoanPaymentChannelsColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.LoanPaymentChannels{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.LoanPaymentChannels{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// DeleteByIDs batch delete loanPaymentChannels by ids
func (d *loanPaymentChannelsDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.LoanPaymentChannels{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.deleteCache(ctx, id)
	}

	return nil
}

// GetByCondition get a loanPaymentChannels by custom condition
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html#_2-condition-parameters-optional
func (d *loanPaymentChannelsDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.LoanPaymentChannels, error) {
	queryStr, args, err := c.ConvertToGorm(query.WithWhitelistNames(model.LoanPaymentChannelsColumnNames))
	if err != nil {
		return nil, err
	}

	table := &model.LoanPaymentChannels{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs Batch get loanPaymentChannels by ids
func (d *loanPaymentChannelsDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanPaymentChannels, error) {
	// no cache
	if d.cache == nil {
		var records []*model.LoanPaymentChannels
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.LoanPaymentChannels)
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
			var records []*model.LoanPaymentChannels
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
				if err = d.cache.MultiSet(ctx, records, cache.LoanPaymentChannelsExpireTime); err != nil {
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

// GetByLastID Get a paginated list of loanPaymentChannelss by last id
func (d *loanPaymentChannelsDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanPaymentChannels, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.LoanPaymentChannels{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *loanPaymentChannelsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanPaymentChannels) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *loanPaymentChannelsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.LoanPaymentChannels{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *loanPaymentChannelsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanPaymentChannels) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
