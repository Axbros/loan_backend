package dao

import (
	"context"
	"errors"
	"loan/internal/types"
	"time"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/go-dev-frame/sponge/pkg/utils"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"loan/internal/cache"
	"loan/internal/database"
	"loan/internal/model"
)

var _ LoanPermissionsDao = (*loanPermissionsDao)(nil)

// LoanPermissionsDao defining the dao interface
type LoanPermissionsDao interface {
	Create(ctx context.Context, table *model.LoanPermissions) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.LoanPermissions) error
	GetByID(ctx context.Context, id uint64) (*model.LoanPermissions, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanPermissions, int64, error)

	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.LoanPermissions, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanPermissions, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanPermissions, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanPermissions) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanPermissions) error

	// 查询 role 当前有效权限 id 列表（deleted_at is null）
	GetActivePermissionIDsByRoleID(ctx context.Context, roleID int64) ([]int64, error)
	// 查询 role 当前已软删权限 id 列表（deleted_at is not null）
	GetDeletedPermissionIDsByRoleID(ctx context.Context, roleID int64) ([]int64, error)

	// 批量恢复（deleted_at = null）
	RestoreByRoleIDAndPermissionIDs(ctx context.Context, roleID int64, permissionIDs []int64) error
	RestoreByTx(ctx context.Context, tx *gorm.DB, roleID int64, permissionIDs []int64) error

	// 批量软删（deleted_at = now）
	SoftDeleteByRoleIDAndPermissionIDs(ctx context.Context, roleID int64, permissionIDs []int64) error
	SoftDeleteByTx(ctx context.Context, tx *gorm.DB, roleID int64, permissionIDs []int64) error

	// 批量创建关联
	BulkCreate(ctx context.Context, roleID int64, permissionIDs []int64) error
	BulkCreateByTx(ctx context.Context, tx *gorm.DB, roleID int64, permissionIDs []int64) error

	// JOIN 查询权限详情（id,name,code）
	GetRolePermissions(ctx context.Context, roleID int64, page, limit int) ([]*types.LoanRolePermissionsObjTable, int64, error)

	SetRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error
}

type loanPermissionsDao struct {
	db    *gorm.DB
	cache cache.LoanPermissionsCache // if nil, the cache is not used.
	sfg   *singleflight.Group        // if cache is nil, the sfg is not used.
}

var (
	ErrInvalidParams      = errors.New("invalid params")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrRoleNotFound       = errors.New("role not found") // 如果你愿意做 role 校验或识别 FK 错误
)

func (d *loanPermissionsDao) SetRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	if roleID <= 0 {
		return ErrInvalidParams
	}

	// 可选：限制一次提交的权限数量，避免恶意/错误请求拖垮 DB
	const maxPermissionIDs = 200
	if len(permissionIDs) > maxPermissionIDs {
		return ErrInvalidParams
	}

	// 1) 去重 + 校验（允许空数组：表示清空）
	targetSet := make(map[int64]struct{}, len(permissionIDs))
	uniq := make([]int64, 0, len(permissionIDs))
	for _, pid := range permissionIDs {
		if pid <= 0 {
			return ErrInvalidParams
		}
		if _, ok := targetSet[pid]; ok {
			continue
		}
		targetSet[pid] = struct{}{}
		uniq = append(uniq, pid)
	}

	// 2) 校验 permission 是否存在（强烈建议）
	if len(uniq) > 0 {
		u64s := make([]uint64, 0, len(uniq))
		for _, id := range uniq {
			u64s = append(u64s, uint64(id))
		}
		m, err := d.GetByIDs(ctx, u64s)
		if err != nil {
			return err
		}
		if len(m) != len(u64s) {
			return ErrPermissionNotFound
		}
	}

	// 3) 事务：差量更新
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 3.1 当前有效权限（用 tx 读）
		activeIDs, err := d.GetActivePermissionIDsByRoleIDTx(ctx, tx, roleID)
		if err != nil {
			return err
		}
		activeSet := make(map[int64]struct{}, len(activeIDs))
		for _, id := range activeIDs {
			activeSet[id] = struct{}{}
		}

		// 3.2 当前已软删权限（用 tx 读）
		deletedIDs, err := d.GetDeletedPermissionIDsByRoleIDTx(ctx, tx, roleID)
		if err != nil {
			return err
		}
		deletedSet := make(map[int64]struct{}, len(deletedIDs))
		for _, id := range deletedIDs {
			deletedSet[id] = struct{}{}
		}

		// 3.3 toRemove = active - target
		toRemove := make([]int64, 0)
		for pid := range activeSet {
			if _, ok := targetSet[pid]; !ok {
				toRemove = append(toRemove, pid)
			}
		}

		// 3.4 toRestore / toInsert = target - active
		toRestore := make([]int64, 0)
		toInsert := make([]int64, 0)
		for pid := range targetSet {
			if _, ok := activeSet[pid]; ok {
				continue
			}
			if _, ok := deletedSet[pid]; ok {
				toRestore = append(toRestore, pid)
			} else {
				toInsert = append(toInsert, pid)
			}
		}

		// 3.5 restore -> insert -> soft delete（都走 tx）
		if err := d.RestoreByTx(ctx, tx, roleID, toRestore); err != nil {
			return err
		}
		if err := d.BulkCreateByTx(ctx, tx, roleID, toInsert); err != nil {
			return err
		}
		if err := d.SoftDeleteByTx(ctx, tx, roleID, toRemove); err != nil {
			return err
		}

		return nil
	})
}

func (d *loanPermissionsDao) GetActivePermissionIDsByRoleIDTx(ctx context.Context, tx *gorm.DB, roleID int64) ([]int64, error) {
	var ids []int64
	err := tx.WithContext(ctx).
		Table("loan_role_permissions").
		Select("permission_id").
		Where("role_id = ? AND deleted_at IS NULL", roleID).
		Scan(&ids).Error
	return ids, err
}

func (d *loanPermissionsDao) GetDeletedPermissionIDsByRoleIDTx(ctx context.Context, tx *gorm.DB, roleID int64) ([]int64, error) {
	var ids []int64
	err := tx.WithContext(ctx).
		Table("loan_role_permissions").
		Select("permission_id").
		Where("role_id = ? AND deleted_at IS NOT NULL", roleID).
		Scan(&ids).Error
	return ids, err
}

// GetActivePermissionIDsByRoleID 查询 role 当前有效权限（deleted_at IS NULL）
func (d *loanPermissionsDao) GetActivePermissionIDsByRoleID(ctx context.Context, roleID int64) ([]int64, error) {
	var ids []int64
	err := d.db.WithContext(ctx).
		Table("loan_role_permissions").
		Select("permission_id").
		Where("role_id = ? AND deleted_at IS NULL", roleID).
		Scan(&ids).Error
	return ids, err
}

// GetDeletedPermissionIDsByRoleID 查询 role 当前已软删权限（deleted_at IS NOT NULL）
func (d *loanPermissionsDao) GetDeletedPermissionIDsByRoleID(ctx context.Context, roleID int64) ([]int64, error) {
	var ids []int64
	err := d.db.WithContext(ctx).
		Table("loan_role_permissions").
		Select("permission_id").
		Where("role_id = ? AND deleted_at IS NOT NULL", roleID).
		Scan(&ids).Error
	return ids, err
}

// RestoreByRoleIDAndPermissionIDs 批量恢复（deleted_at = NULL）
func (d *loanPermissionsDao) RestoreByRoleIDAndPermissionIDs(ctx context.Context, roleID int64, permissionIDs []int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return d.RestoreByTx(ctx, tx, roleID, permissionIDs)
	})
}

// RestoreByTx 批量恢复（deleted_at = NULL）- 使用外部事务
func (d *loanPermissionsDao) RestoreByTx(ctx context.Context, tx *gorm.DB, roleID int64, permissionIDs []int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}
	now := time.Now()
	return tx.WithContext(ctx).
		Table("loan_role_permissions").
		Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).
		Updates(map[string]any{
			"deleted_at": nil,
			"updated_at": now,
		}).Error
}

// SoftDeleteByRoleIDAndPermissionIDs 批量软删（deleted_at = NOW）
func (d *loanPermissionsDao) SoftDeleteByRoleIDAndPermissionIDs(ctx context.Context, roleID int64, permissionIDs []int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return d.SoftDeleteByTx(ctx, tx, roleID, permissionIDs)
	})
}

// SoftDeleteByTx 批量软删（deleted_at = NOW）- 使用外部事务
func (d *loanPermissionsDao) SoftDeleteByTx(ctx context.Context, tx *gorm.DB, roleID int64, permissionIDs []int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}
	now := time.Now()
	return tx.WithContext(ctx).
		Table("loan_role_permissions").
		Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).
		Updates(map[string]any{
			"deleted_at": now,
			"updated_at": now,
		}).Error
}

// BulkCreate 批量创建 role-permission 关联（忽略重复主键）
func (d *loanPermissionsDao) BulkCreate(ctx context.Context, roleID int64, permissionIDs []int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return d.BulkCreateByTx(ctx, tx, roleID, permissionIDs)
	})
}

// BulkCreateByTx 批量创建 - 使用外部事务
func (d *loanPermissionsDao) BulkCreateByTx(ctx context.Context, tx *gorm.DB, roleID int64, permissionIDs []int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	now := time.Now()
	rows := make([]map[string]interface{}, 0, len(permissionIDs))
	for _, pid := range permissionIDs {
		rows = append(rows, map[string]interface{}{
			"role_id":       roleID,
			"permission_id": pid,
			"created_at":    now,
			"updated_at":    now,
			"deleted_at":    nil,
		})
	}

	return tx.WithContext(ctx).
		Table("loan_role_permissions").
		Clauses(clause.Insert{Modifier: "IGNORE"}). // ✅ MySQL: INSERT IGNORE
		Create(&rows).Error
}

// GetRolePermissions JOIN 查询角色权限详情（id,name,code），只返回未删除的关联
func (d *loanPermissionsDao) GetRolePermissions(ctx context.Context, roleID int64, page, limit int) ([]*types.LoanRolePermissionsObjTable, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if page < 0 {
		page = 0
	}
	offset := page * limit

	base := d.db.WithContext(ctx).
		Table("loan_role_permissions AS r").
		Joins("INNER JOIN loan_permissions p ON r.permission_id = p.id").
		Where("r.role_id = ? AND r.deleted_at IS NULL", roleID)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	records := make([]*types.LoanRolePermissionsObjTable, 0, limit)
	err := base.
		Select("r.id AS id, p.name AS name, p.code AS code").
		Order("r.id ASC").
		Limit(limit).
		Offset(offset).
		Scan(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// NewLoanPermissionsDao creating the dao interface
func NewLoanPermissionsDao(db *gorm.DB, xCache cache.LoanPermissionsCache) LoanPermissionsDao {
	if xCache == nil {
		return &loanPermissionsDao{db: db}
	}
	return &loanPermissionsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *loanPermissionsDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new loanPermissions, insert the record and the id value is written back to the table
func (d *loanPermissionsDao) Create(ctx context.Context, table *model.LoanPermissions) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a loanPermissions by id
func (d *loanPermissionsDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.LoanPermissions{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a loanPermissions by ids
func (d *loanPermissionsDao) UpdateByID(ctx context.Context, table *model.LoanPermissions) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *loanPermissionsDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.LoanPermissions) error {
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
	if table.Type != "" {
		update["type"] = table.Type
	}
	if table.Resource != "" {
		update["resource"] = table.Resource
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a loanPermissions by id
func (d *loanPermissionsDao) GetByID(ctx context.Context, id uint64) (*model.LoanPermissions, error) {
	// no cache
	if d.cache == nil {
		record := &model.LoanPermissions{}
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
			table := &model.LoanPermissions{}
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
			if err = d.cache.Set(ctx, id, table, cache.LoanPermissionsExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.LoanPermissions)
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

// GetByColumns get a paginated list of loanPermissionss by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *loanPermissionsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanPermissions, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.LoanPermissionsColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.LoanPermissions{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.LoanPermissions{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// DeleteByIDs batch delete loanPermissions by ids
func (d *loanPermissionsDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.LoanPermissions{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.deleteCache(ctx, id)
	}

	return nil
}

// GetByCondition get a loanPermissions by custom condition
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html#_2-condition-parameters-optional
func (d *loanPermissionsDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.LoanPermissions, error) {
	queryStr, args, err := c.ConvertToGorm(query.WithWhitelistNames(model.LoanPermissionsColumnNames))
	if err != nil {
		return nil, err
	}

	table := &model.LoanPermissions{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs Batch get loanPermissions by ids
func (d *loanPermissionsDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanPermissions, error) {
	// no cache
	if d.cache == nil {
		var records []*model.LoanPermissions
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.LoanPermissions)
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
			var records []*model.LoanPermissions
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
				if err = d.cache.MultiSet(ctx, records, cache.LoanPermissionsExpireTime); err != nil {
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

// GetByLastID Get a paginated list of loanPermissionss by last id
func (d *loanPermissionsDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanPermissions, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.LoanPermissions{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *loanPermissionsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanPermissions) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *loanPermissionsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.LoanPermissions{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *loanPermissionsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanPermissions) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
