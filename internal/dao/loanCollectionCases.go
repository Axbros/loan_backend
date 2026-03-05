package dao

import (
	"context"
	"errors"
	"loan/internal/types"
	"strings"
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

var _ LoanCollectionCasesDao = (*loanCollectionCasesDao)(nil)

// LoanCollectionCasesDao defining the dao interface
type LoanCollectionCasesDao interface {
	Create(ctx context.Context, table *model.LoanCollectionCases) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.LoanCollectionCases) error
	GetByID(ctx context.Context, id uint64) (*model.LoanCollectionCases, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*types.LoanCollectionCasesObjTable, int64, error)

	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.LoanCollectionCases, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanCollectionCases, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanCollectionCases, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanCollectionCases) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanCollectionCases) error
}

type loanCollectionCasesDao struct {
	db    *gorm.DB
	cache cache.LoanCollectionCasesCache // if nil, the cache is not used.
	sfg   *singleflight.Group            // if cache is nil, the sfg is not used.
}

// NewLoanCollectionCasesDao creating the dao interface
func NewLoanCollectionCasesDao(db *gorm.DB, xCache cache.LoanCollectionCasesCache) LoanCollectionCasesDao {
	if xCache == nil {
		return &loanCollectionCasesDao{db: db}
	}
	return &loanCollectionCasesDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *loanCollectionCasesDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new loanCollectionCases, insert the record and the id value is written back to the table
func (d *loanCollectionCasesDao) Create(ctx context.Context, table *model.LoanCollectionCases) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a loanCollectionCases by id
func (d *loanCollectionCasesDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.LoanCollectionCases{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a loanCollectionCases by ids
func (d *loanCollectionCasesDao) UpdateByID(ctx context.Context, table *model.LoanCollectionCases) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *loanCollectionCasesDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.LoanCollectionCases) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.ScheduleID != 0 {
		update["schedule_id"] = table.ScheduleID
	}
	if table.CollectorUserID != 0 {
		update["collector_user_id"] = table.CollectorUserID
	}
	if table.AssignedByUserID != 0 {
		update["assigned_by_user_id"] = table.AssignedByUserID
	}

	if table.Priority != 0 {
		update["priority"] = table.Priority
	}
	if table.Status != 0 {
		update["status"] = table.Status
	}

	if table.CompletedAt != nil && table.CompletedAt.IsZero() == false {
		update["completed_at"] = table.CompletedAt
	}
	if table.CompletedNote != "" {
		update["completed_note"] = table.CompletedNote
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a loanCollectionCases by id
func (d *loanCollectionCasesDao) GetByID(ctx context.Context, id uint64) (*model.LoanCollectionCases, error) {
	// no cache
	if d.cache == nil {
		record := &model.LoanCollectionCases{}
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
			table := &model.LoanCollectionCases{}
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
			if err = d.cache.Set(ctx, id, table, cache.LoanCollectionCasesExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.LoanCollectionCases)
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

// GetByColumns get a paginated list of loanCollectionCasess by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *loanCollectionCasesDao) GetByColumns(ctx context.Context, params *query.Params) ([]*types.LoanCollectionCasesObjTable, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(
		query.WithWhitelistNames(model.LoanCollectionCasesColumnNames),
	)
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	// base query（等价你的 SQL）
	base := d.db.WithContext(ctx).
		Table("loan_collection_cases AS cc").
		Joins("INNER JOIN loan_repayment_schedules AS rs ON cc.schedule_id = rs.id").
		Joins("INNER JOIN loan_users AS collector ON cc.collector_user_id = collector.id").
		Joins("LEFT JOIN loan_users AS assigner ON cc.assigned_by_user_id = assigner.id").
		Joins("INNER JOIN loan_disbursements d ON d.id = rs.disbursement_id").
		Joins("INNER JOIN loan_baseinfo b ON b.id = d.baseinfo_id").
		Where("cc.deleted_at IS NULL").
		Where("rs.deleted_at IS NULL").
		Where("d.deleted_at IS NULL").
		Where("b.deleted_at IS NULL")

	// apply dynamic conditions（注意：你的白名单字段要能匹配 cc.*，否则需要在 ConvertToGormConditions 前缀列名）
	if queryStr != "" {
		base = base.Where(prefixCC(queryStr), args...)
	}

	// count（联表可能放大行数，保险起见用 distinct cc.id）
	var total int64
	if params.Sort != "ignore count" {
		if err := base.Distinct("cc.id").Count(&total).Error; err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, 0, nil
		}
	}

	// page/sort
	order, limit, offset := params.ConvertToPage()

	records := make([]*types.LoanCollectionCasesObjTable, 0, limit)
	err = base.
		Select(`
			cc.id,
			cc.schedule_id,
			d.baseinfo_id,
			b.first_name,
			b.second_name,
			b.age,
			b.gender,
			b.id_type,
			b.id_number,
			b.mobile,
			cc.priority,
			cc.status,
			cc.completed_at,
			cc.completed_note,
			rs.due_date,
			d.net_amount,
			rs.total_due,
			rs.paid_total,
			cc.created_at,
			collector.username AS collector_name,
			assigner.username AS assigned_by_name
		`).
		Order(order).
		Limit(limit).
		Offset(offset).
		Scan(&records).Error

	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// DeleteByIDs batch delete loanCollectionCases by ids
func (d *loanCollectionCasesDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.LoanCollectionCases{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.deleteCache(ctx, id)
	}

	return nil
}

// GetByCondition get a loanCollectionCases by custom condition
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html#_2-condition-parameters-optional
func (d *loanCollectionCasesDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.LoanCollectionCases, error) {
	queryStr, args, err := c.ConvertToGorm(query.WithWhitelistNames(model.LoanCollectionCasesColumnNames))
	if err != nil {
		return nil, err
	}

	table := &model.LoanCollectionCases{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs Batch get loanCollectionCases by ids
func (d *loanCollectionCasesDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanCollectionCases, error) {
	// no cache
	if d.cache == nil {
		var records []*model.LoanCollectionCases
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.LoanCollectionCases)
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
			var records []*model.LoanCollectionCases
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
				if err = d.cache.MultiSet(ctx, records, cache.LoanCollectionCasesExpireTime); err != nil {
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

// GetByLastID Get a paginated list of loanCollectionCasess by last id
func (d *loanCollectionCasesDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanCollectionCases, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.LoanCollectionCases{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *loanCollectionCasesDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanCollectionCases) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *loanCollectionCasesDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.LoanCollectionCases{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *loanCollectionCasesDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanCollectionCases) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

// 给 LoanCollectionCases 的列名都加 cc. 前缀，避免歧义
func prefixCC(queryStr string) string {
	// 按你们 whitelist 里常用列逐个替换（你可以把需要的都加上）
	replaces := map[string]string{
		" status ":              " cc.status ",
		" priority ":            " cc.priority ",
		" schedule_id ":         " cc.schedule_id ",
		" collector_user_id ":   " cc.collector_user_id ",
		" assigned_by_user_id ": " cc.assigned_by_user_id ",
		" created_at ":          " cc.created_at ",
		" updated_at ":          " cc.updated_at ",
		" deleted_at ":          " cc.deleted_at ",
	}

	// 处理开头/结尾等情况（=、<、>、IN 等）
	// 这里用更通用的写法：替换 "status" 为 "cc.status"（但要避免替换到别名/字符串里）
	// 最简单可用版本（你们 queryStr 规则一般很规整）：按关键 token 替换
	for k, v := range replaces {
		queryStr = strings.ReplaceAll(queryStr, k, v)
	}
	// 兼容 queryStr 以 "status" 开头
	queryStr = strings.ReplaceAll(queryStr, "status =", "cc.status =")
	queryStr = strings.ReplaceAll(queryStr, "status IN", "cc.status IN")
	queryStr = strings.ReplaceAll(queryStr, "status !=", "cc.status !=")
	queryStr = strings.ReplaceAll(queryStr, "status >", "cc.status >")
	queryStr = strings.ReplaceAll(queryStr, "status <", "cc.status <")

	return queryStr
}
