package dao

import (
	"context"
	"errors"
	"fmt"
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
	GetByColumnsWithAuditRecords(ctx context.Context, params *query.Params) ([]*model.LoanBaseinfoWithAuditRecord, int64, error)
	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.LoanBaseinfo, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanBaseinfo, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanBaseinfo, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanBaseinfo) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanBaseinfo) error

	GetFilesMapByBaseinfoID(ctx context.Context, baseinfoID uint64) (map[string][]string, error)
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

func (d *loanBaseinfoDao) GetFilesMapByBaseinfoID(ctx context.Context, baseinfoID uint64) (map[string][]string, error) {
	type row struct {
		Type   string `gorm:"column:type"`
		OssURL string `gorm:"column:oss_url"`
	}

	var rows []row
	err := d.db.WithContext(ctx).
		Table("loan_baseinfo_files").
		Select("type, oss_url").
		Where("baseinfo_id = ? AND deleted_at IS NULL", baseinfoID).
		Order("type ASC, id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	res := make(map[string][]string)
	for _, r := range rows {
		res[r.Type] = append(res[r.Type], r.OssURL)
	}
	return res, nil
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

	if table.HasCar != 0 {
		update["has_car"] = table.HasCar
	}
	if table.Mobile != "" {
		update["mobile"] = table.Mobile
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

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a loanBaseinfo by id
func (d *loanBaseinfoDao) GetByID(ctx context.Context, id uint64) (*model.LoanBaseinfo, error) {
	// no cache：无缓存直接查库，查完填充风险信息
	if d.cache == nil {
		record := &model.LoanBaseinfo{}
		err := d.db.WithContext(ctx).Where("id = ?", id).First(record).Error
		if err != nil {
			return record, err
		}
		// 核心新增：为单条数据填充风险信息（复用批量填充方法，兼容单条）
		if fillErr := d.fillRiskCustomerInfo(ctx, []*model.LoanBaseinfo{record}); fillErr != nil {
			logger.Warn("fillRiskCustomerInfo error", logger.Err(fillErr), logger.Any("loan_baseinfo_id", id))
		}
		return record, nil
	}

	// get from cache：优先从缓存获取（缓存中已包含风险字段，直接返回）
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	// get from database：缓存未命中，查库并做高并发加锁、防穿透
	if errors.Is(err, database.ErrCacheNotFound) {
		// 相同ID加锁，防止高并发同时查库（击穿数据库）
		val, err, _ := d.sfg.Do(utils.Uint64ToStr(id), func() (interface{}, error) {
			table := &model.LoanBaseinfo{}
			err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
			if err != nil {
				// 数据库无记录，设置缓存占位符防穿透
				if errors.Is(err, database.ErrRecordNotFound) {
					if setErr := d.cache.SetPlaceholder(ctx, id); setErr != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(setErr), logger.Any("id", id))
					}
					return nil, database.ErrRecordNotFound
				}
				return nil, err
			}

			// 核心新增：查库成功后，填充风险信息
			if fillErr := d.fillRiskCustomerInfo(ctx, []*model.LoanBaseinfo{table}); fillErr != nil {
				logger.Warn("fillRiskCustomerInfo error", logger.Err(fillErr), logger.Any("loan_baseinfo_id", id))
			}

			// 缓存回写（此时table已包含风险字段，缓存存储完整数据）
			if setErr := d.cache.Set(ctx, id, table, cache.LoanBaseinfoExpireTime); setErr != nil {
				logger.Warn("cache.Set error", logger.Err(setErr), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		// 类型断言，转换为LoanBaseinfo指针
		table, ok := val.(*model.LoanBaseinfo)
		if !ok {
			return nil, database.ErrRecordNotFound
		}
		return table, nil
	}

	// 缓存返回占位符错误，说明数据库无此记录，返回记录不存在
	if d.cache.IsPlaceholderErr(err) {
		return nil, database.ErrRecordNotFound
	}

	// 其他缓存错误，直接返回
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

// 定义 LoanUser 模型（如果项目中已有，无需重复定义，确保字段匹配）
// type LoanUser struct {
// 	ID       uint64 `gorm:"column:id;primaryKey" json:"id"`
// 	Username string `gorm:"column:username;type:varchar(64)" json:"username"` // 真实姓名
// 	// 其他字段（如手机号、创建时间等）按需保留
// }

func (d *loanBaseinfoDao) GetByColumnsWithAuditRecords(ctx context.Context, params *query.Params) ([]*model.LoanBaseinfoWithAuditRecord, int64, error) {
	// 1. 转换查询参数
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.LoanBaseinfoColumnNames))
	if err != nil {
		return nil, 0, fmt.Errorf("query params error: %w", err)
	}

	var total int64
	const ignoreCountSortFlag = "ignore count"
	if params.Sort != ignoreCountSortFlag {
		err = d.db.WithContext(ctx).Model(&model.LoanBaseinfo{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, fmt.Errorf("count loan baseinfo error: %w", err)
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	// 2. 查询贷款基础信息
	records := []*model.LoanBaseinfo{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).
		Order(order).
		Limit(limit).
		Offset(offset).
		Where(queryStr, args...).
		Find(&records).Error
	if err != nil {
		return nil, 0, fmt.Errorf("query loan baseinfo list error: %w", err)
	}

	// 3. 批量查询审核记录
	var baseinfoIDs []uint64
	for _, record := range records {
		if record.ID != 0 {
			baseinfoIDs = append(baseinfoIDs, record.ID)
		}
	}

	auditRecordMap := make(map[uint64][]*model.LoanAudits)
	var allAuditorUserIDs []uint64 // 收集所有审核人员ID
	if len(baseinfoIDs) > 0 {
		var auditRecords []*model.LoanAudits
		err = d.db.WithContext(ctx).
			Model(&model.LoanAudits{}).
			Where("baseinfo_id IN (?) and audit_result = 1", baseinfoIDs). //只提取审核成功的
			Find(&auditRecords).Error
		if err != nil {
			return nil, 0, fmt.Errorf("batch query audit records error: %w", err)
		}

		// 构建审核记录映射表 + 收集审核人员ID
		for _, ar := range auditRecords {
			auditRecordMap[ar.BaseinfoID] = append(auditRecordMap[ar.BaseinfoID], ar)
			// 收集非空的审核人员ID
			if ar.AuditorUserID != 0 {
				allAuditorUserIDs = append(allAuditorUserIDs, ar.AuditorUserID)
			}
		}
	}

	// ========== 核心新增：批量查询审核人员姓名 ==========
	// 4. 批量查询 loan_users 表，获取审核人员姓名
	auditorNameMap := make(map[uint64]string) // key: auditorUserID, value: username
	if len(allAuditorUserIDs) > 0 {
		var loanUsers []*model.LoanUsers
		err = d.db.WithContext(ctx).
			Model(&model.LoanUsers{}).
			Select("id, username"). // 只查需要的字段，提升性能
			Where("id IN (?)", allAuditorUserIDs).
			Find(&loanUsers).Error
		if err != nil {
			return nil, 0, fmt.Errorf("batch query loan users error: %w", err)
		}

		// 构建 ID->姓名 映射表
		for _, user := range loanUsers {
			auditorNameMap[user.ID] = user.Username
		}
	}

	// 5. 转换为结果结构体 + 替换审核人员姓名
	results := make([]*model.LoanBaseinfoWithAuditRecord, 0, len(records))
	for _, record := range records {
		if auditRecordMap[record.ID] == nil {
			//表示审核拒绝的
			continue
		}
		result := &model.LoanBaseinfoWithAuditRecord{
			LoanBaseinfo: *record,
			AuditRecords: auditRecordMap[record.ID],
		}

		// 遍历审核记录，补充 auditorName 字段
		for _, auditRecord := range result.AuditRecords {
			// 从映射表获取姓名，无则显示"未知"
			auditRecord.AuditorName = auditorNameMap[auditRecord.AuditorUserID]
			if auditRecord.AuditorName == "" {
				auditRecord.AuditorName = "未知"
			}
		}

		results = append(results, result)
	}

	return results, total, nil
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

func (d *loanBaseinfoDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanBaseinfo, error) {
	// no cache：无缓存时直接查库+关联风险
	if d.cache == nil {
		var records []*model.LoanBaseinfo
		// 1. 查询贷款基础信息主表
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		// 2. 批量关联风险客户表数据（核心新增）
		if err = d.fillRiskCustomerInfo(ctx, records); err != nil {
			logger.Warn("fillRiskCustomerInfo error", logger.Err(err), logger.Any("ids", ids))
			// 仅打警告，不阻断主数据返回
		}
		// 3. 转map返回
		itemMap := make(map[uint64]*model.LoanBaseinfo)
		for _, record := range records {
			itemMap[record.ID] = record
		}
		return itemMap, nil
	}

	// get from cache：优先从缓存获取（缓存中已包含风险字段）
	itemMap, err := d.cache.MultiGet(ctx, ids)
	if err != nil {
		return nil, err
	}

	// 筛选缓存未命中的ID
	var missedIDs []uint64
	for _, id := range ids {
		if _, ok := itemMap[id]; !ok {
			missedIDs = append(missedIDs, id)
		}
	}

	// 处理缓存未命中的ID：查库+关联风险+缓存回写+占位符
	if len(missedIDs) > 0 {
		var realMissedIDs []uint64
		for _, id := range missedIDs {
			_, err = d.cache.Get(ctx, id)
			if d.cache.IsPlaceholderErr(err) {
				continue
			}
			realMissedIDs = append(realMissedIDs, id)
		}

		if len(realMissedIDs) > 0 {
			var records []*model.LoanBaseinfo
			var recordIDMap = make(map[uint64]struct{})
			// 1. 查数据库主表
			err = d.db.WithContext(ctx).Where("id IN (?)", realMissedIDs).Find(&records).Error
			if err != nil {
				return nil, err
			}
			if len(records) > 0 {
				// 2. 核心新增：为数据库查到的主表数据，批量关联风险信息
				if err = d.fillRiskCustomerInfo(ctx, records); err != nil {
					logger.Warn("fillRiskCustomerInfo error", logger.Err(err), logger.Any("ids", realMissedIDs))
				}
				// 3. 更新结果map+标记已查到的ID
				for _, record := range records {
					itemMap[record.ID] = record
					recordIDMap[record.ID] = struct{}{}
				}
				// 4. 缓存回写（此时records已包含风险字段，缓存会存储完整数据）
				if err = d.cache.MultiSet(ctx, records, cache.LoanBaseinfoExpireTime); err != nil {
					logger.Warn("cache.MultiSet error", logger.Err(err), logger.Any("ids", records))
				}
				if len(records) == len(realMissedIDs) {
					return itemMap, nil
				}
			}
			// 5. 为数据库不存在的ID设置缓存占位符
			for _, id := range realMissedIDs {
				if _, ok := recordIDMap[id]; !ok {
					if err = d.cache.SetPlaceholder(ctx, id); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("id", id))
					}
				}
			}
		}
	}

	// 返回包含风险信息的结果map
	return itemMap, nil
}

// 核心新增方法：批量为LoanBaseinfo填充风险客户表的信息
// 入参是贷款基础信息切片，批量查询风险数据后赋值，避免循环单查提升性能
// ========== 核心改造：fillRiskCustomerInfo 方法（新增created_by关联username逻辑）==========
// 功能升级：批量填充风险信息 + 批量关联操作人ID(created_by)查询用户名，仅2次数据库查询，无N+1
func (d *loanBaseinfoDao) fillRiskCustomerInfo(ctx context.Context, baseinfos []*model.LoanBaseinfo) error {
	if len(baseinfos) == 0 {
		return nil
	}

	// 步骤1：提取贷款基础信息ID，批量查询风险表（含created_by操作人ID）
	var baseinfoIDs []uint64
	for _, b := range baseinfos {
		baseinfoIDs = append(baseinfoIDs, b.ID)
	}
	var riskList []*model.LoanRiskCustomer
	err := d.db.WithContext(ctx).Where("loan_baseinfo_id IN (?)", baseinfoIDs).Find(&riskList).Error
	if err != nil {
		return err
	}
	// 构建「贷款基础信息ID→风险记录」的映射，方便后续赋值
	riskMap := make(map[uint64]*model.LoanRiskCustomer)
	// 收集所有风险记录的操作人ID（created_by），用于批量查询用户表
	var operateUserIDs []uint64
	for _, r := range riskList {
		riskMap[r.LoanBaseinfoID] = r
		if r.CreatedBy > 0 { // 过滤无效操作人ID（0）
			operateUserIDs = append(operateUserIDs, r.CreatedBy)
		}
	}

	// 步骤2：批量查询操作人信息（loan_users），仅查需要的ID和username，提升效率
	operateUserMap := make(map[uint64]string) // 「操作人ID→用户名」的映射
	if len(operateUserIDs) > 0 {
		var users []*model.LoanUsers
		err = d.db.WithContext(ctx).
			Select("id, username"). // 仅查必要字段，减少数据传输
			Where("id IN (?)", operateUserIDs).
			Find(&users).Error
		if err != nil {
			logger.Warn("query loan_users for risk operate error", logger.Err(err), logger.Any("operateUserIDs", operateUserIDs))
			// 仅打警告，不阻断后续赋值（操作人用户名置空即可）
		} else {
			for _, u := range users {
				operateUserMap[uint64(u.ID)] = u.Username // 注意类型转换：user.id是int64，created_by是uint64
			}
		}
	}

	// 步骤3：批量赋值风险信息 + 操作人信息到LoanBaseinfo
	for _, b := range baseinfos {
		risk, ok := riskMap[b.ID]
		if !ok {
			// 无风险记录：重置所有风险/操作人字段为默认值
			b.RiskListStatus = 0
			b.RiskListReason = ""
			b.RiskListMarkedAt = nil
			b.RiskOperateID = 0
			b.RiskOperateName = ""
			continue
		}

		// 赋值原有风险信息
		switch risk.RiskType {
		case -1:
			b.RiskListStatus = 2 // 黑名单
		case 1:
			b.RiskListStatus = 1 // 白名单
		default:
			b.RiskListStatus = 0 // 正常
		}
		b.RiskListReason = risk.RiskReason
		b.RiskListMarkedAt = &risk.CreatedAt // 风险标记时间=风险记录创建时间

		// 核心新增：赋值操作人信息（created_by关联）
		b.RiskOperateID = risk.CreatedBy                   // 操作人ID（loan_risk_customer.created_by）
		b.RiskOperateName = operateUserMap[risk.CreatedBy] // 操作人用户名（关联loan_users.username）
	}

	return nil
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
