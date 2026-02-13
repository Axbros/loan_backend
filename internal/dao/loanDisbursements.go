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

var _ LoanDisbursementsDao = (*loanDisbursementsDao)(nil)

// LoanDisbursementsDao defining the dao interface
type LoanDisbursementsDao interface {
	Create(ctx context.Context, table *model.LoanDisbursements) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.LoanDisbursements) error
	GetByID(ctx context.Context, id uint64) (*model.LoanDisbursements, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanDisbursements, int64, error)

	DeleteByIDs(ctx context.Context, ids []uint64) error
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.LoanDisbursements, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanDisbursements, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanDisbursements, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanDisbursements) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanDisbursements) error

	GetOverviewList(ctx context.Context, req *types.ListLoanDisbursementsOverviewRequest) (*types.ListLoanDisbursementsOverviewResponse, error)
}

type loanDisbursementsDao struct {
	db    *gorm.DB
	cache cache.LoanDisbursementsCache // if nil, the cache is not used.
	sfg   *singleflight.Group          // if cache is nil, the sfg is not used.
}

// NewLoanDisbursementsDao creating the dao interface
func NewLoanDisbursementsDao(db *gorm.DB, xCache cache.LoanDisbursementsCache) LoanDisbursementsDao {
	if xCache == nil {
		return &loanDisbursementsDao{db: db}
	}
	return &loanDisbursementsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *loanDisbursementsDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new loanDisbursements, insert the record and the id value is written back to the table
func (d *loanDisbursementsDao) Create(ctx context.Context, table *model.LoanDisbursements) error {
	return d.db.WithContext(ctx).Create(table).Error
}

func (d *loanDisbursementsDao) GetOverviewList(
	ctx context.Context,
	req *types.ListLoanDisbursementsOverviewRequest) (*types.ListLoanDisbursementsOverviewResponse, error) {
	// 1. 定义返回结果
	var (
		list  []*types.LoanDisbursedList
		total int64

		// 存储动态条件和参数
		whereConditions []string
		whereArgs       []interface{}
	)

	// -------------------------- 拼接 loan_baseinfo 表的过滤条件 --------------------------
	if req.Condition != nil {
		// 1. 姓名：模糊查询
		if req.Condition.Name != "" {
			whereConditions = append(whereConditions, "b.first_name LIKE ?")
			whereArgs = append(whereArgs, "%"+req.Condition.Name+"%")
		}

		// 2. 年龄：精确匹配
		if req.Condition.Age != nil && *req.Condition.Age > 0 {
			whereConditions = append(whereConditions, "b.age = ?")
			whereArgs = append(whereArgs, *req.Condition.Age)
		}

		// 3. 性别：精确匹配（M/W）
		if req.Condition.Gender != "" {
			whereConditions = append(whereConditions, "b.gender = ?")
			whereArgs = append(whereArgs, req.Condition.Gender)
		}

		// 4. 证件类型：精确匹配
		if req.Condition.IDType != "" {
			whereConditions = append(whereConditions, "b.id_type = ?")
			whereArgs = append(whereArgs, req.Condition.IDType)
		}

		// 5. 证件号码：精确匹配
		if req.Condition.IDNo != "" {
			whereConditions = append(whereConditions, "b.id_number = ?")
			whereArgs = append(whereArgs, req.Condition.IDNo)
		}

		// 6. 申请金额：精确匹配
		if req.Condition.LoanAmount != nil && *req.Condition.LoanAmount > 0 {
			whereConditions = append(whereConditions, "b.application_amount = ?")
			whereArgs = append(whereArgs, *req.Condition.LoanAmount)
		}
	}

	// -------------------------- 构造 SQL --------------------------
	baseSQL := `
        FROM
            loan_disbursements d
            INNER JOIN loan_baseinfo b ON d.baseinfo_id = b.id
            INNER JOIN loan_payment_channels c ON d.payout_channel_id = c.id
    `

	// 拼接 WHERE 子句
	whereSQL := ""
	if len(whereConditions) > 0 {
		whereSQL = " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// 1. 查询总条数
	countSQL := "SELECT COUNT(*) " + baseSQL + whereSQL
	err := d.db.WithContext(ctx).Raw(countSQL, whereArgs...).Scan(&total).Error
	if err != nil {
		logger.Error("查询放款概览总条数失败", logger.Err(err), logger.Any("req.Condition", req.Condition))
		return nil, err
	}

	// 2. 计算分页偏移量
	offset := req.Page * req.Limit

	// 3. 分页查询列表
	querySQL := `
        SELECT
            b.id,
            b.first_name,
            b.age,
            b.gender,
            b.id_type,
            b.id_number,
            b.application_amount,
            d.net_amount,
            b.loan_days,
            c.name,
            d.payout_order_no,
            c.payout_fee_rate
    ` + baseSQL + whereSQL + " LIMIT ? OFFSET ?"

	// 合并分页参数
	queryArgs := append(whereArgs, req.Limit, offset)

	// 执行查询
	err = d.db.WithContext(ctx).Raw(querySQL, queryArgs...).Scan(&list).Error
	if err != nil {
		logger.Error("分页查询放款概览列表失败", logger.Err(err), logger.Any("pageReq", req), logger.Any("filter", req.Condition))
		return nil, err
	}

	// -------------------------- 返回结果 --------------------------
	return &types.ListLoanDisbursementsOverviewResponse{
		Total: total,
		List:  list,
	}, nil
}

// DeleteByID delete a loanDisbursements by id
func (d *loanDisbursementsDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.LoanDisbursements{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a loanDisbursements by ids
func (d *loanDisbursementsDao) UpdateByID(ctx context.Context, table *model.LoanDisbursements) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *loanDisbursementsDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.LoanDisbursements) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.BaseinfoID != 0 {
		update["baseinfo_id"] = table.BaseinfoID
	}
	if table.DisburseAmount != 0 {
		update["disburse_amount"] = table.DisburseAmount
	}
	if table.NetAmount != 0 {
		update["net_amount"] = table.NetAmount
	}
	if table.Status != 0 {
		update["status"] = table.Status
	}
	if table.SourceReferrerUserID != 0 {
		update["source_referrer_user_id"] = table.SourceReferrerUserID
	}
	if table.AuditorUserID != 0 {
		update["auditor_user_id"] = table.AuditorUserID
	}
	if table.AuditedAt != nil && table.AuditedAt.IsZero() == false {
		update["audited_at"] = table.AuditedAt
	}
	if table.PayoutChannelID != 0 {
		update["payout_channel_id"] = table.PayoutChannelID
	}
	if table.PayoutOrderNo != "" {
		update["payout_order_no"] = table.PayoutOrderNo
	}
	if table.DisbursedAt != nil && table.DisbursedAt.IsZero() == false {
		update["disbursed_at"] = table.DisbursedAt
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a loanDisbursements by id
func (d *loanDisbursementsDao) GetByID(ctx context.Context, id uint64) (*model.LoanDisbursements, error) {
	// no cache
	if d.cache == nil {
		record := &model.LoanDisbursements{}
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
			table := &model.LoanDisbursements{}
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
			if err = d.cache.Set(ctx, id, table, cache.LoanDisbursementsExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.LoanDisbursements)
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

// GetByColumns get a paginated list of loanDisbursementss by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *loanDisbursementsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.LoanDisbursements, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.LoanDisbursementsColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.LoanDisbursements{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.LoanDisbursements{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// DeleteByIDs batch delete loanDisbursements by ids
func (d *loanDisbursementsDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.LoanDisbursements{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.deleteCache(ctx, id)
	}

	return nil
}

// GetByCondition get a loanDisbursements by custom condition
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html#_2-condition-parameters-optional
func (d *loanDisbursementsDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.LoanDisbursements, error) {
	queryStr, args, err := c.ConvertToGorm(query.WithWhitelistNames(model.LoanDisbursementsColumnNames))
	if err != nil {
		return nil, err
	}

	table := &model.LoanDisbursements{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs Batch get loanDisbursements by ids
func (d *loanDisbursementsDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.LoanDisbursements, error) {
	// no cache
	if d.cache == nil {
		var records []*model.LoanDisbursements
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.LoanDisbursements)
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
			var records []*model.LoanDisbursements
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
				if err = d.cache.MultiSet(ctx, records, cache.LoanDisbursementsExpireTime); err != nil {
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

// GetByLastID Get a paginated list of loanDisbursementss by last id
func (d *loanDisbursementsDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.LoanDisbursements, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.LoanDisbursements{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Limit()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *loanDisbursementsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanDisbursements) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *loanDisbursementsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.LoanDisbursements{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *loanDisbursementsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanDisbursements) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
