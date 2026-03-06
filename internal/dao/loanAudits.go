package dao

import (
	"context"
	"errors"
	"loan/internal/types"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"loan/internal/cache"
	"loan/internal/model"

	"github.com/go-dev-frame/sponge/pkg/logger"
)

var _ LoanAuditsDao = (*loanAuditsDao)(nil)

// LoanAuditsDao defining the dao interface
type LoanAuditsDao interface {
	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanAudits) (uint64, error)
	ListByBaseinfoID(ctx context.Context, baseinfoID uint64) ([]*model.LoanAudits, error)
	GetByBaseinfoID(ctx context.Context, baseinfoID uint64, auditType int) (*types.LoanAuditDetail, error)
	GetDisbursmentsByBaseInfoID(ctx context.Context, baseInfoID uint64) (*types.DisbursementWithChannel, error)
}

type loanAuditsDao struct {
	db    *gorm.DB
	cache cache.LoanAuditsCache // if nil, the cache is not used.
	sfg   *singleflight.Group   // if cache is nil, the sfg is not used.
}

// NewLoanAuditsDao creating the dao interface
func NewLoanAuditsDao(db *gorm.DB, xCache cache.LoanAuditsCache) LoanAuditsDao {
	if xCache == nil {
		return &loanAuditsDao{db: db}
	}
	return &loanAuditsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *loanAuditsDao) GetDisbursmentsByBaseInfoID(ctx context.Context, baseInfoID uint64) (*types.DisbursementWithChannel, error) {
	var result types.DisbursementWithChannel
	// 执行你需要的关联查询SQL
	sqlStr := `
        SELECT 
            d.id,
            d.disburse_amount,
            d.net_amount,
            d.status,
            d.payout_order_no,
            d.disbursed_at,
            c.name AS channel_name
        FROM 
            loan_disbursements d
        INNER JOIN 
            loan_payment_channels c ON d.payout_channel_id = c.id
        WHERE 
            d.baseinfo_id = ?;
    `

	// 执行原生SQL并扫描结果
	err := d.db.WithContext(ctx).Raw(sqlStr, baseInfoID).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("查询放款记录及渠道信息失败",
				logger.String("reason", "未找到记录"),
				logger.Uint64("baseinfo_id", baseInfoID),
			)
			return nil, nil // 无数据返回nil，上层处理
		}
		logger.Error("关联查询放款记录+渠道失败", logger.Err(err))
		return nil, err
	}

	return &result, nil
}

func (d *loanAuditsDao) GetByBaseinfoID(ctx context.Context, baseinfoID uint64, auditType int) (*types.LoanAuditDetail, error) {
	// 关键修改1：初始化具体的结构体实例，而非 nil 指针
	auditRecord := &types.LoanAuditDetail{}
	// 核心查询部分
	err := d.db.WithContext(ctx).Model(&model.LoanAudits{}).
		Select("a.*, u.username as auditor_username").
		Joins("INNER JOIN loan_users u ON a.auditor_user_id = u.id").
		Where("a.baseinfo_id =? and a.audit_type = ?", baseinfoID, auditType).
		Table("loan_audits a").  // 给 loan_audits 起别名 a
		First(auditRecord).Error // 关键修改2：用 First 替代 Find，查询单条记录

	// 关键修改3：处理 "记录未找到" 的场景
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 查询成功但无数据，返回 nil 和 nil
			return nil, nil
		}
		// 其他查询错误，返回 nil 和具体错误
		return nil, err
	}

	// 有数据时返回结构体指针和 nil 错误
	return auditRecord, nil
}

func (d *loanAuditsDao) ListByBaseinfoID(ctx context.Context, baseinfoID uint64) ([]*model.LoanAudits, error) {
	var records []*model.LoanAudits
	err := d.db.WithContext(ctx).
		Where("baseinfo_id = ? AND deleted_at IS NULL", baseinfoID).
		Order("created_at DESC, id DESC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *loanAuditsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.LoanAudits) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}
