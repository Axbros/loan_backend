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

var _ LoanRepaymentSchedulesDao = (*loanRepaymentSchedulesDao)(nil)

// LoanRepaymentSchedulesDao defining the dao interface
type LoanRepaymentSchedulesDao interface {
	Create(ctx context.Context, table *model.LoanRepaymentSchedules) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.LoanRepaymentSchedules) error
	GetByID(ctx context.Context, id uint64) (*model.LoanRepaymentSchedules, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanRepaymentSchedules, int64, error)

	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.LoanRepaymentSchedules, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanRepaymentSchedules, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanRepaymentSchedules, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanRepaymentSchedules) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanRepaymentSchedules) error
}

type loanRepaymentSchedulesDao struct {
	db    *gorm.DB
	cache cache.LoanRepaymentSchedulesCache // if nil, the cache is not used.
	sfg   *singleflight.Group               // if cache is nil, the sfg is not used.
}

// NewLoanRepaymentSchedulesDao creating the dao interface
func NewLoanRepaymentSchedulesDao(db *gorm.DB, xCache cache.LoanRepaymentSchedulesCache) LoanRepaymentSchedulesDao {
	if xCache == nil {
		return &loanRepaymentSchedulesDao{db: db}
	}
	return &loanRepaymentSchedulesDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *loanRepaymentSchedulesDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new loanRepaymentSchedules, insert the record and the id value is written back to the table
func (d *loanRepaymentSchedulesDao) Create(ctx context.Context, table *model.LoanRepaymentSchedules) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a loanRepaymentSchedules by id
func (d *loanRepaymentSchedulesDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.LoanRepaymentSchedules{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a loanRepaymentSchedules by ids
func (d *loanRepaymentSchedulesDao) UpdateByID(ctx context.Context, table *model.LoanRepaymentSchedules) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *loanRepaymentSchedulesDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.LoanRepaymentSchedules) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.DisbursementID != 0 {
		update["disbursement_id"] = table.DisbursementID
	}
	if table.InstallmentNo != 0 {
		update["installment_no"] = table.InstallmentNo
	}
	if table.DueDate != nil && table.DueDate.IsZero() == false {
		update["due_date"] = table.DueDate
	}
	if table.PrincipalDue != 0 {
		update["principal_due"] = table.PrincipalDue
	}
	if table.InterestDue != 0 {
		update["interest_due"] = table.InterestDue
	}
	if table.FeeDue != 0 {
		update["fee_due"] = table.FeeDue
	}
	if table.PenaltyDue != 0 {
		update["penalty_due"] = table.PenaltyDue
	}
	if table.TotalDue != 0 {
		update["total_due"] = table.TotalDue
	}
	if table.PaidPrincipal != 0 {
		update["paid_principal"] = table.PaidPrincipal
	}
	if table.PaidInterest != 0 {
		update["paid_interest"] = table.PaidInterest
	}
	if table.PaidFee != 0 {
		update["paid_fee"] = table.PaidFee
	}
	if table.PaidPenalty != 0 {
		update["paid_penalty"] = table.PaidPenalty
	}
	if table.PaidTotal != 0 {
		update["paid_total"] = table.PaidTotal
	}
	if table.Status != 0 {
		update["status"] = table.Status
	}
	if table.LastPaidAt != nil && table.LastPaidAt.IsZero() == false {
		update["last_paid_at"] = table.LastPaidAt
	}
	if table.SettledAt != nil && table.SettledAt.IsZero() == false {
		update["settled_at"] = table.SettledAt
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a loanRepaymentSchedules by id
func (d *loanRepaymentSchedulesDao) GetByID(ctx context.Context, id uint64) (*model.LoanRepaymentSchedules, error) {
	// no cache
	if d.cache == nil {
		record := &model.LoanRepaymentSchedules{}
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
			table := &model.LoanRepaymentSchedules{}
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
			if err = d.cache.Set(ctx, id, table, cache.LoanRepaymentSchedulesExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.LoanRepaymentSchedules)
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

// GetByColumns get a paginated list of loanRepaymentScheduless by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *loanRepaymentSchedulesDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanRepaymentSchedules, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.LoanRepaymentSchedulesColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.LoanRepaymentSchedules{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.LoanRepaymentSchedules{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// DeleteByIDs batch delete loanRepaymentSchedules by ids
func (d *loanRepaymentSchedulesDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.LoanRepaymentSchedules{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.deleteCache(ctx, id)
	}

	return nil
}

// GetByCondition get a loanRepaymentSchedules by custom condition
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html#_2-condition-parameters-optional
func (d *loanRepaymentSchedulesDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.LoanRepaymentSchedules, error) {
	queryStr, args, err := c.ConvertToGorm(query.WithWhitelistNames(model.LoanRepaymentSchedulesColumnNames))
	if err != nil {
		return nil, err
	}

	table := &model.LoanRepaymentSchedules{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs Batch get loanRepaymentSchedules by ids
func (d *loanRepaymentSchedulesDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanRepaymentSchedules, error) {
	// no cache
	if d.cache == nil {
		var records []*model.LoanRepaymentSchedules
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.LoanRepaymentSchedules)
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
			var records []*model.LoanRepaymentSchedules
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
				if err = d.cache.MultiSet(ctx, records, cache.LoanRepaymentSchedulesExpireTime); err != nil {
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

// GetByLastID Get a paginated list of loanRepaymentScheduless by last id
func (d *loanRepaymentSchedulesDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanRepaymentSchedules, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.LoanRepaymentSchedules{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *loanRepaymentSchedulesDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanRepaymentSchedules) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *loanRepaymentSchedulesDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.LoanRepaymentSchedules{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *loanRepaymentSchedulesDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanRepaymentSchedules) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
