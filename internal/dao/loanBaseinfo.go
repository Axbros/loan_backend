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

var _ LoanBaseinfoDao = (*loanBaseinfoDao)(nil)

// LoanBaseinfoDao defining the dao interface
type LoanBaseinfoDao interface {
	Create(ctx context.Context, table *model.LoanBaseinfo) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.LoanBaseinfo) error
	GetByID(ctx context.Context, id uint64) (*model.LoanBaseinfo, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanBaseinfo, int64, error)

	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.LoanBaseinfo, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanBaseinfo, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanBaseinfo, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanBaseinfo) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanBaseinfo) error
}

type loanBaseinfoDao struct {
	db    *gorm.DB
	cache cache.LoanBaseinfoCache // if nil, the cache is not used.
	sfg   *singleflight.Group     // if cache is nil, the sfg is not used.
}

// NewLoanBaseinfoDao creating the dao interface
func NewLoanBaseinfoDao(db *gorm.DB, xCache cache.LoanBaseinfoCache) LoanBaseinfoDao {
	if xCache == nil {
		return &loanBaseinfoDao{db: db}
	}
	return &loanBaseinfoDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *loanBaseinfoDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new loanBaseinfo, insert the record and the id value is written back to the table
func (d *loanBaseinfoDao) Create(ctx context.Context, table *model.LoanBaseinfo) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a loanBaseinfo by id
func (d *loanBaseinfoDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.LoanBaseinfo{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a loanBaseinfo by ids
func (d *loanBaseinfoDao) UpdateByID(ctx context.Context, table *model.LoanBaseinfo) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *loanBaseinfoDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.LoanBaseinfo) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.FirstName != "" {
		update["first_name"] = table.FirstName
	}
	if table.SecondName != "" {
		update["second_name"] = table.SecondName
	}
	if table.Age != 0 {
		update["age"] = table.Age
	}
	if table.Gender != "" {
		update["gender"] = table.Gender
	}
	if table.IdType != "" {
		update["id_type"] = table.IdType
	}
	if table.IdNumber != "" {
		update["id_number"] = table.IdNumber
	}
	if table.IdCard != "" {
		update["id_card"] = table.IdCard
	}
	if table.Operator != "" {
		update["operator"] = table.Operator
	}
	if table.Inviter != "" {
		update["inviter"] = table.Inviter
	}
	if table.Work != "" {
		update["work"] = table.Work
	}
	if table.Company != "" {
		update["company"] = table.Company
	}
	if table.Salary != 0 {
		update["salary"] = table.Salary
	}
	if table.MaritalStatus != 0 {
		update["marital_status"] = table.MaritalStatus
	}
	if table.HasHouse != 0 {
		update["has_house"] = table.HasHouse
	}
	if table.PropertyCertificate != "" {
		update["property_certificate"] = table.PropertyCertificate
	}
	if table.HasCar != 0 {
		update["has_car"] = table.HasCar
	}
	if table.VehicleRgistrationCertificate != "" {
		update["vehicle_rgistration_certificate"] = table.VehicleRgistrationCertificate
	}
	if table.ApplicationAmount != 0 {
		update["application_amount"] = table.ApplicationAmount
	}
	if table.AuditStatus != 0 {
		update["audit_status"] = table.AuditStatus
	}
	if table.BankNo != "" {
		update["bank_no"] = table.BankNo
	}
	if table.ClientIP != "" {
		update["client_ip"] = table.ClientIP
	}
	if table.ReferrerUserID != 0 {
		update["referrer_user_id"] = table.ReferrerUserID
	}
	if table.RefCode != "" {
		update["ref_code"] = table.RefCode
	}
	if table.LoanDays != 0 {
		update["loan_days"] = table.LoanDays
	}
	if table.RiskListStatus != 0 {
		update["risk_list_status"] = table.RiskListStatus
	}
	if table.RiskListReason != "" {
		update["risk_list_reason"] = table.RiskListReason
	}
	if table.RiskListMarkedAt != nil && table.RiskListMarkedAt.IsZero() == false {
		update["risk_list_marked_at"] = table.RiskListMarkedAt
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a loanBaseinfo by id
func (d *loanBaseinfoDao) GetByID(ctx context.Context, id uint64) (*model.LoanBaseinfo, error) {
	// no cache
	if d.cache == nil {
		record := &model.LoanBaseinfo{}
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
			table := &model.LoanBaseinfo{}
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
			if err = d.cache.Set(ctx, id, table, cache.LoanBaseinfoExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.LoanBaseinfo)
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

// GetByColumns get a paginated list of loanBaseinfos by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *loanBaseinfoDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanBaseinfo, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.LoanBaseinfoColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.LoanBaseinfo{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.LoanBaseinfo{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// DeleteByIDs batch delete loanBaseinfo by ids
func (d *loanBaseinfoDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.LoanBaseinfo{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.deleteCache(ctx, id)
	}

	return nil
}

// GetByCondition get a loanBaseinfo by custom condition
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html#_2-condition-parameters-optional
func (d *loanBaseinfoDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.LoanBaseinfo, error) {
	queryStr, args, err := c.ConvertToGorm(query.WithWhitelistNames(model.LoanBaseinfoColumnNames))
	if err != nil {
		return nil, err
	}

	table := &model.LoanBaseinfo{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs Batch get loanBaseinfo by ids
func (d *loanBaseinfoDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanBaseinfo, error) {
	// no cache
	if d.cache == nil {
		var records []*model.LoanBaseinfo
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.LoanBaseinfo)
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
			var records []*model.LoanBaseinfo
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
				if err = d.cache.MultiSet(ctx, records, cache.LoanBaseinfoExpireTime); err != nil {
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

// GetByLastID Get a paginated list of loanBaseinfos by last id
func (d *loanBaseinfoDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanBaseinfo, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.LoanBaseinfo{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *loanBaseinfoDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanBaseinfo) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *loanBaseinfoDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.LoanBaseinfo{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *loanBaseinfoDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanBaseinfo) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
