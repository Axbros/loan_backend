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

var _ LoanUsersDao = (*loanUsersDao)(nil)

// LoanUsersDao defining the dao interface
type LoanUsersDao interface {
	Create(ctx context.Context, table *model.LoanUsers) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.LoanUsers) error
	GetByID(ctx context.Context, id uint64) (*model.LoanUsers, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanUsers, int64, error)

	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.LoanUsers, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanUsers, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanUsers, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanUsers) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanUsers) error

	GetByUsername(ctx context.Context, username string) (*model.LoanUsers, error)
	GetRoleCodesByUserID(ctx context.Context, uid uint64) ([]string, error)
	GetPermCodesByUserID(ctx context.Context, uid uint64) ([]string, error)
}

type loanUsersDao struct {
	db    *gorm.DB
	cache cache.LoanUsersCache // if nil, the cache is not used.
	sfg   *singleflight.Group  // if cache is nil, the sfg is not used.
}

// NewLoanUsersDao creating the dao interface
func NewLoanUsersDao(db *gorm.DB, xCache cache.LoanUsersCache) LoanUsersDao {
	if xCache == nil {
		return &loanUsersDao{db: db}
	}
	return &loanUsersDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}
func (d *loanUsersDao) GetRoleCodesByUserID(ctx context.Context, userID uint64) ([]string, error) {
	var roleCodes []string

	err := d.db.WithContext(ctx).
		Table("loan_user_roles ur").
		Select("r.code").
		Joins("JOIN loan_roles r ON r.id = ur.role_id AND r.deleted_at IS NULL").
		Where("ur.user_id = ? AND ur.deleted_at IS NULL", userID).
		Where("r.status = 1").
		Order("r.code ASC").
		Scan(&roleCodes).Error
	if err != nil {
		return nil, err
	}

	if roleCodes == nil {
		roleCodes = []string{}
	}
	return roleCodes, nil
}
func (d *loanUsersDao) GetPermCodesByUserID(ctx context.Context, userID uint64) ([]string, error) {
	var permCodes []string

	err := d.db.WithContext(ctx).
		Table("loan_user_roles ur").
		Select("DISTINCT p.code").
		Joins("JOIN loan_roles r ON r.id = ur.role_id AND r.deleted_at IS NULL AND r.status = 1").
		Joins("JOIN loan_role_permissions rp ON rp.role_id = ur.role_id AND rp.deleted_at IS NULL").
		Joins("JOIN loan_permissions p ON p.id = rp.permission_id AND p.deleted_at IS NULL").
		Where("ur.user_id = ? AND ur.deleted_at IS NULL", userID).
		Order("p.code ASC").
		Scan(&permCodes).Error
	if err != nil {
		return nil, err
	}

	if permCodes == nil {
		permCodes = []string{}
	}
	return permCodes, nil
}

func (d *loanUsersDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new loanUsers, insert the record and the id value is written back to the table
func (d *loanUsersDao) Create(ctx context.Context, table *model.LoanUsers) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a loanUsers by id
func (d *loanUsersDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.LoanUsers{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a loanUsers by ids
func (d *loanUsersDao) UpdateByID(ctx context.Context, table *model.LoanUsers) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *loanUsersDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.LoanUsers) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.Username != "" {
		update["username"] = table.Username
	}
	if table.PasswordHash != "" {
		update["password_hash"] = table.PasswordHash
	}
	if table.DepartmentID != 0 {
		update["department_id"] = table.DepartmentID
	}
	if table.MfaEnabled != 0 {
		update["mfa_enabled"] = table.MfaEnabled
	}
	if table.MfaRequired != 0 {
		update["mfa_required"] = table.MfaRequired
	}
	if table.Status != 0 {
		update["status"] = table.Status
	}
	if table.ShareCode != "" {
		update["share_code"] = table.ShareCode
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a loanUsers by id
func (d *loanUsersDao) GetByID(ctx context.Context, id uint64) (*model.LoanUsers, error) {
	// no cache
	if d.cache == nil {
		record := &model.LoanUsers{}
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
			table := &model.LoanUsers{}
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
			if err = d.cache.Set(ctx, id, table, cache.LoanUsersExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.LoanUsers)
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

// GetByColumns get a paginated list of loanUserss by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *loanUsersDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanUsers, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.LoanUsersColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.LoanUsers{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.LoanUsers{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// DeleteByIDs batch delete loanUsers by ids
func (d *loanUsersDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.LoanUsers{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.deleteCache(ctx, id)
	}

	return nil
}

// GetByCondition get a loanUsers by custom condition
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html#_2-condition-parameters-optional
func (d *loanUsersDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.LoanUsers, error) {
	queryStr, args, err := c.ConvertToGorm(query.WithWhitelistNames(model.LoanUsersColumnNames))
	if err != nil {
		return nil, err
	}

	table := &model.LoanUsers{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs Batch get loanUsers by ids
func (d *loanUsersDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanUsers, error) {
	// no cache
	if d.cache == nil {
		var records []*model.LoanUsers
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.LoanUsers)
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
			var records []*model.LoanUsers
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
				if err = d.cache.MultiSet(ctx, records, cache.LoanUsersExpireTime); err != nil {
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

// GetByLastID Get a paginated list of loanUserss by last id
func (d *loanUsersDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanUsers, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.LoanUsers{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (d *loanUsersDao) GetByUsername(ctx context.Context, username string) (*model.LoanUsers, error) {
	record := &model.LoanUsers{}

	err := d.db.WithContext(ctx).
		Where("username = ? AND deleted_at IS NULL", username).
		First(record).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return record, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *loanUsersDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanUsers) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *loanUsersDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.LoanUsers{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *loanUsersDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanUsers) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
