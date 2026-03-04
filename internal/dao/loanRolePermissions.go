package dao

import (
	"context"
	"errors"
	"loan/internal/types"
	"regexp"
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

var _ LoanRolePermissionsDao = (*loanRolePermissionsDao)(nil)

// LoanRolePermissionsDao defining the dao interface
type LoanRolePermissionsDao interface {
	Create(ctx context.Context, table *model.LoanRolePermissions) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.LoanRolePermissions) error
	GetByID(ctx context.Context, id uint64) (*model.LoanRolePermissions, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*types.LoanRolePermissionsObjTable, int64, error)

	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.LoanRolePermissions, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanRolePermissions, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanRolePermissions, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanRolePermissions) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanRolePermissions) error
}

type loanRolePermissionsDao struct {
	db    *gorm.DB
	cache cache.LoanRolePermissionsCache // if nil, the cache is not used.
	sfg   *singleflight.Group            // if cache is nil, the sfg is not used.
}

// NewLoanRolePermissionsDao creating the dao interface
func NewLoanRolePermissionsDao(db *gorm.DB, xCache cache.LoanRolePermissionsCache) LoanRolePermissionsDao {
	if xCache == nil {
		return &loanRolePermissionsDao{db: db}
	}
	return &loanRolePermissionsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *loanRolePermissionsDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new loanRolePermissions, insert the record and the id value is written back to the table
func (d *loanRolePermissionsDao) Create(ctx context.Context, table *model.LoanRolePermissions) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a loanRolePermissions by id
func (d *loanRolePermissionsDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.LoanRolePermissions{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a loanRolePermissions by ids
func (d *loanRolePermissionsDao) UpdateByID(ctx context.Context, table *model.LoanRolePermissions) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *loanRolePermissionsDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.LoanRolePermissions) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.RoleID != 0 {
		update["role_id"] = table.RoleID
	}
	if table.PermissionID != 0 {
		update["permission_id"] = table.PermissionID
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a loanRolePermissions by id
func (d *loanRolePermissionsDao) GetByID(ctx context.Context, id uint64) (*model.LoanRolePermissions, error) {
	// no cache
	if d.cache == nil {
		record := &model.LoanRolePermissions{}
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
			table := &model.LoanRolePermissions{}
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
			if err = d.cache.Set(ctx, id, table, cache.LoanRolePermissionsExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.LoanRolePermissions)
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

func (d *loanRolePermissionsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*types.LoanRolePermissionsObjTable, int64, error) {
	// 1) 生成 where 条件（白名单仍用 LoanRolePermissions 表的列）
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.LoanRolePermissionsColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	// 2) 把 where 里的列名加上 r. 前缀（避免 join 后列名歧义）
	//    例：role_id = ?  ->  r.role_id = ?
	queryStr = prefixWhereColumns(queryStr, model.LoanRolePermissionsColumnNames, "r")

	// 3) 构造基础 join 查询（用于 count 和 list）
	base := d.db.WithContext(ctx).
		Table("loan_role_permissions AS r").
		Joins("INNER JOIN loan_permissions p ON r.permission_id = p.id").
		Where(queryStr, args...).Where("r.deleted_at IS NULL")

	// 4) count
	var total int64
	if params.Sort != "ignore count" {
		if err := base.Count(&total).Error; err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, 0, nil
		}
	}

	// 5) list
	records := make([]*types.LoanRolePermissionsObjTable, 0)
	order, limit, offset := params.ConvertToPage()

	// 可选：如果你的 order 里是 "id desc" 这种，建议也加前缀
	// 否则可能出现 “column id is ambiguous”
	order = prefixOrder(order)

	err = base.
		Select("r.id AS id, p.name AS name, p.code AS code").
		Order(order).
		Limit(limit).
		Offset(offset).
		Scan(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// prefixWhereColumns 把 where 里的列名替换成 alias.col（只替换白名单列）
func prefixWhereColumns(where string, cols map[string]bool, alias string) string {
	for c := range cols {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(c) + `\b`)
		where = re.ReplaceAllString(where, alias+"."+c)
	}
	return where
}

// prefixOrder 给 order 加前缀（简单处理常见场景）
// 你如果允许按 name/code 排序，这里也可以扩展映射到 p.name / p.code
func prefixOrder(order string) string {
	o := strings.TrimSpace(order)
	if o == "" {
		return o
	}

	// 常见：id desc / id asc
	// 让它变成 r.id desc
	if strings.HasPrefix(o, "id ") || o == "id" {
		return "r." + o
	}
	// 如果你的 params 允许 name/code 排序，可以加：
	if strings.HasPrefix(o, "name ") || o == "name" {
		return "p." + o
	}
	if strings.HasPrefix(o, "code ") || o == "code" {
		return "p." + o
	}

	return o
}

// DeleteByIDs batch delete loanRolePermissions by ids
func (d *loanRolePermissionsDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.LoanRolePermissions{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.deleteCache(ctx, id)
	}

	return nil
}

// GetByCondition get a loanRolePermissions by custom condition
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html#_2-condition-parameters-optional
func (d *loanRolePermissionsDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.LoanRolePermissions, error) {
	queryStr, args, err := c.ConvertToGorm(query.WithWhitelistNames(model.LoanRolePermissionsColumnNames))
	if err != nil {
		return nil, err
	}

	table := &model.LoanRolePermissions{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs Batch get loanRolePermissions by ids
func (d *loanRolePermissionsDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanRolePermissions, error) {
	// no cache
	if d.cache == nil {
		var records []*model.LoanRolePermissions
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.LoanRolePermissions)
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
			var records []*model.LoanRolePermissions
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
				if err = d.cache.MultiSet(ctx, records, cache.LoanRolePermissionsExpireTime); err != nil {
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

// GetByLastID Get a paginated list of loanRolePermissionss by last id
func (d *loanRolePermissionsDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanRolePermissions, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.LoanRolePermissions{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *loanRolePermissionsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanRolePermissions) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *loanRolePermissionsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.LoanRolePermissions{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *loanRolePermissionsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanRolePermissions) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
