package handler

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"loan/internal/cache"
	"loan/internal/dao"
	"loan/internal/database"
	"loan/internal/ecode"
	"loan/internal/model"
	"loan/internal/types"
)

var _ LoanBaseinfoHandler = (*loanBaseinfoHandler)(nil)

// LoanBaseinfoHandler defining the handler interface
type LoanBaseinfoHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Review(c *gin.Context)
	WithAuditRecordList(c *gin.Context)

	DeleteByIDs(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	ListByLastID(c *gin.Context)
}

type loanBaseinfoHandler struct {
	iDao                 dao.LoanBaseinfoDao
	auditDao             dao.LoanAuditsDao
	userDao              dao.LoanUsersDao
	channelDao           dao.LoanPaymentChannelsDao
	disbursmentDao       dao.LoanDisbursementsDao
	repaymentScheduleDao dao.LoanRepaymentSchedulesDao
}

// NewLoanBaseinfoHandler creating the handler interface
func NewLoanBaseinfoHandler() LoanBaseinfoHandler {
	return &loanBaseinfoHandler{
		iDao: dao.NewLoanBaseinfoDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanBaseinfoCache(database.GetCacheType()),
		),
		auditDao: dao.NewLoanAuditsDao(
			database.GetDB(),
			cache.NewLoanAuditsCache(database.GetCacheType()),
		),
		userDao: dao.NewLoanUsersDao(
			database.GetDB(),
			cache.NewLoanUsersCache(database.GetCacheType()),
		),
		channelDao: dao.NewLoanPaymentChannelsDao(
			database.GetDB(),
			cache.NewLoanPaymentChannelsCache(database.GetCacheType()),
		),
		disbursmentDao: dao.NewLoanDisbursementsDao(
			database.GetDB(),
			cache.NewLoanDisbursementsCache(database.GetCacheType()),
		),
		repaymentScheduleDao: dao.NewLoanRepaymentSchedulesDao(
			database.GetDB(),
			cache.NewLoanRepaymentSchedulesCache(database.GetCacheType()),
		),
	}
}

type Audit_Status int

const (
	Rejected      Audit_Status = -1 // 机审拒绝
	Pending       Audit_Status = 0
	PreReview     Audit_Status = 1 //初审通过，等待财务审核
	FinanceReview Audit_Status = 2 //财务审核通过，最终台
)

type AuditType int // 修正原Audit_Type命名，符合Go大驼峰规范
const (
	PreReviewType     = 1 //初审审核
	FinanceReviewType = 2 //放款审核
	IncomeReviewType  = 3 //回款审核
)

func (h *loanBaseinfoHandler) WithAuditRecordList(c *gin.Context) {
	form := &types.ListLoanBaseinfosRequestWithAuditType{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	if len(form.Columns) == 0 {
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)

	loanBaseinfos, total, err := h.iDao.GetByColumnsWithAuditRecords(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertSimpleLoanBaseinfosWithAuditRecord(loanBaseinfos)
	if err != nil {
		response.Error(c, ecode.ErrListLoanBaseinfo)
		return
	}

	response.Success(c, gin.H{
		"records": data,
		"total":   total,
	})
}

// 新增：将请求的audit_type字符串（0/1/2）转换为AuditType枚举，同时做合法性校验
func parseAuditType(auditTypeStr int) AuditType {
	switch auditTypeStr {
	case 1:
		return PreReviewType
	case 2:
		return FinanceReviewType
	case 3:
		return IncomeReviewType
	default:
		// 返回自定义的"无效审核类型"错误码（需在ecode中定义）
		return -1
	}
}

func (h *loanBaseinfoHandler) Review(c *gin.Context) {
	ctx := middleware.WrapCtx(c)

	uid, ok := getUIDFromClaims(c)
	if !ok || uid == 0 {
		response.Out(c, ecode.Unauthorized)
		return
	}

	form := &types.AuditRequest{}
	if err := c.ShouldBindJSON(form); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	auditType := parseAuditType(form.AuditType)
	if auditType == -1 {
		logger.Warn("无效的审核类型", logger.String("audit_type", strconv.Itoa(form.AuditType)), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidAuditType)
		return
	}

	// 0) MFA 必须启用（不进事务）
	u, err := h.userDao.GetByID(ctx, uid)
	if err != nil {
		logger.Error("查询用户信息失败", logger.Err(err), logger.Uint64("uid", uid), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGetByIDLoanUsers)
		return
	}
	if u.MfaEnabled != 1 {
		logger.Warn("用户未启用MFA，禁止审核", logger.Uint64("uid", uid), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.MFANotEnabled)
		return
	}

	// 1) MFA 校验（不进事务）
	otpCode := strings.TrimSpace(form.MfaCode)
	if len(otpCode) != 6 {
		logger.Warn("MFA验证码长度错误", logger.String("otp_code", otpCode), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.MFAOTPRequired)
		return
	}

	dev, err := h.userDao.GetActivePrimaryMFADevice(ctx, uid)
	if err != nil {
		logger.Error("查询MFA主设备失败", logger.Err(err), logger.Uint64("uid", uid), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InternalServerError)
		return
	}
	if dev == nil {
		logger.Warn("用户无激活的MFA主设备", logger.Uint64("uid", uid), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGetByIDLoanMfaDevices)
		return
	}

	secret, err := decryptSecretFromBytes(dev.SecretEnc)
	if err != nil {
		logger.Error("解密MFA密钥失败", logger.Err(err), logger.Uint64("device_id", dev.ID), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrSecret)
		return
	}
	secret = strings.TrimSpace(secret)
	if secret == "" {
		logger.Warn("MFA密钥为空", logger.Uint64("device_id", dev.ID), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrSecret)
		return
	}

	ok, err = totp.ValidateCustom(otpCode, secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		logger.Error("校验MFA验证码失败", logger.Err(err), logger.Uint64("uid", uid), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrSecret)
		return
	}
	if !ok {
		logger.Warn("MFA验证码无效", logger.String("otp_code", otpCode), logger.Uint64("uid", uid), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidOTP)
		return
	}

	// 2) 开事务（银行级：显式 begin/commit/rollback + recover）
	db := database.GetDB()
	tx := db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback().Error
			logger.Error("transaction panic", logger.Any("panic", r), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.InternalServerError)
		}
	}()

	if err := tx.Error; err != nil {
		logger.Error("tx begin error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InternalServerError)
		return
	}

	// 3) 事务内 DAO（推荐 WithTx；没有的话看文末替代方案）
	iDao := dao.NewLoanBaseinfoDao(tx, cache.NewLoanBaseinfoCache(database.GetCacheType()))
	auditDao := dao.NewLoanAuditsDao(tx, cache.NewLoanAuditsCache(database.GetCacheType()))
	channelDao := dao.NewLoanPaymentChannelsDao(tx, cache.NewLoanPaymentChannelsCache(database.GetCacheType()))
	disDao := dao.NewLoanDisbursementsDao(tx, cache.NewLoanDisbursementsCache(database.GetCacheType()))
	repayDao := dao.NewLoanRepaymentSchedulesDao(tx, cache.NewLoanRepaymentSchedulesCache(database.GetCacheType()))

	// 4) 银行级：锁住 baseinfo 行，防止并发双审/双放款
	loanBaseinfoRecord := &model.LoanBaseinfo{}
	if err := tx.
		// 需要：import "gorm.io/gorm/clause"
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", form.CustomerID).
		First(loanBaseinfoRecord).Error; err != nil {

		_ = tx.Rollback().Error
		logger.Warn("baseinfo not found / lock failed", logger.Err(err), logger.Uint64("customer_id", form.CustomerID), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGetByIDLoanBaseinfo)
		return
	}

	// 5) 状态机校验（按你们业务可调整）
	//    示例规则：
	//    - 已拒绝（-1）不允许再审核
	//    - 财务审核必须在前置审核通过后才能执行（如果你们有这种要求）
	if loanBaseinfoRecord.AuditStatus == -1 {
		_ = tx.Rollback().Error
		response.Error(c, ecode.InvalidParams) // 你可以换成更合适的 ecode，例如 ErrAlreadyRejected
		return
	}

	// 如果你们要求按顺序：
	// if form.AuditType == FinanceReviewType && loanBaseinfoRecord.AuditStatus != PreReviewType {
	//     _ = tx.Rollback().Error
	//     response.Error(c, ecode.InvalidParams) // 换成 ErrInvalidAuditFlow
	//     return
	// }

	// 6) 财务审核通过：幂等防重复放款（强烈建议）
	//    做法：对 disbursement 做 FOR UPDATE 查询，若已存在则不再创建（避免重复放款）
	var createdDisbursementID uint64
	if form.AuditType == FinanceReviewType && form.AuditResult {
		if form.PaymentChannelID == 0 {
			_ = tx.Rollback().Error
			response.Error(c, ecode.ErrPaymentChannel)
			return
		}

		// 可选：对支付渠道也加锁（一般不需要）
		paymentChannelRecord, err := channelDao.GetByID(ctx, form.PaymentChannelID)
		if err != nil {
			_ = tx.Rollback().Error
			response.Error(c, ecode.ErrGetByIDLoanPaymentChannels)
			return
		}
		if paymentChannelRecord == nil {
			_ = tx.Rollback().Error
			response.Error(c, ecode.ErrGetByIDLoanPaymentChannels)
			return
		}

		// 幂等检查：是否已经有 disbursement（用 tx 查并锁）
		// 你需要有一个能按 baseinfo_id 查 disbursement 的查询；
		// 如果你没有 dao 方法，就直接用 tx.Where 查表（下面示例直接用 tx）
		existing := &model.LoanDisbursements{}
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("baseinfo_id = ?", form.CustomerID).
			Order("id DESC").
			First(existing).Error

		if err == nil && existing != nil && existing.ID != 0 {
			// ✅ 已放款（或已创建放款记录），这里按“幂等成功”处理：
			// - 不再创建 disbursement/schedule
			createdDisbursementID = existing.ID
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// ✅ 不存在：创建 disbursement + schedule
			if loanBaseinfoRecord.ApplicationAmount == 0 {
				_ = tx.Rollback().Error
				response.Error(c, ecode.InvalidParams)
				return
			}

			var feeRate float32 = 0.0
			if paymentChannelRecord.PayoutFeeRate != 0 {
				feeRate = paymentChannelRecord.PayoutFeeRate // 如 35 → 代表 35%
			}
			// 1. 校验申请金额（int64类型，分）
			if loanBaseinfoRecord.ApplicationAmount == 0 {
				_ = tx.Rollback().Error
				response.Error(c, ecode.InvalidParams)
				return
			}

			applicationAmount := loanBaseinfoRecord.ApplicationAmount
			// 先转成int64计算，避免float32精度丢失（核心！）
			feeAmount := (applicationAmount * int64(feeRate)) / 100 // 手续费（分）
			netAmount := applicationAmount - feeAmount              // 净金额（分）
			disburseAmount := applicationAmount                     //单位是分

			now := time.Now()
			currentTime := &now

			disbursmentRecord := &model.LoanDisbursements{
				BaseinfoID:           form.CustomerID,
				DisburseAmount:       disburseAmount,
				NetAmount:            netAmount,
				Status:               1, // 已放款（如你们需要先“放款中”，可改为中间态）
				SourceReferrerUserID: loanBaseinfoRecord.ReferrerUserID,
				AuditorUserID:        uid,
				PayoutChannelID:      form.PaymentChannelID,
				AuditedAt:            currentTime,
				DisbursedAt:          currentTime,
				PayoutOrderNo:        generateOrderNo(), // 建议数据库层加唯一索引保证不重复
			}

			if err := disDao.Create(ctx, disbursmentRecord); err != nil {
				_ = tx.Rollback().Error
				response.Error(c, ecode.ErrCreateLoanDisbursements)
				return
			}
			createdDisbursementID = disbursmentRecord.ID

			// 还款计划：建议也做幂等（对 disbursement_id 唯一）
			loanDays := loanBaseinfoRecord.LoanDays
			if loanDays <= 0 {
				loanDays = 30
			}
			dueDate := now.AddDate(0, 0, loanDays)
			dueDateOnly := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, dueDate.Location())
			dueDatePtr := &dueDateOnly

			scheduleRecord := &model.LoanRepaymentSchedules{
				DisbursementID: int64(createdDisbursementID),
				InstallmentNo:  1,
				DueDate:        dueDatePtr,
				PrincipalDue:   disburseAmount,
				InterestDue:    0,
				FeeDue:         disburseAmount,
				PenaltyDue:     0,
				TotalDue:       disburseAmount,
				PaidPrincipal:  0,
				PaidInterest:   0,
				PaidFee:        0,
				PaidPenalty:    0,
				PaidTotal:      0,
				Status:         0,
			}

			if err := repayDao.Create(ctx, scheduleRecord); err != nil {
				_ = tx.Rollback().Error
				response.Error(c, ecode.ErrCreateLoanRepaymentSchedules)
				return
			}
		} else {
			// 查询出错
			_ = tx.Rollback().Error
			response.Error(c, ecode.InternalServerError)
			return
		}
	}

	// 7) 更新 baseinfo 审核状态（通过/拒绝都执行）
	auditResult := -1
	if form.AuditResult {
		auditResult = 1
		loanBaseinfoRecord.AuditStatus = form.AuditType
	} else {
		loanBaseinfoRecord.AuditStatus = -1
	}

	if err := iDao.UpdateByID(ctx, loanBaseinfoRecord); err != nil {
		_ = tx.Rollback().Error
		response.Error(c, ecode.ErrUpdateByIDLoanBaseinfo)
		return
	}

	// 8) 创建审核记录（通过/拒绝都执行）
	record := &model.LoanAudits{
		AuditResult:   auditResult,
		AuditType:     int(auditType),
		AuditComment:  form.Remark,
		BaseinfoID:    form.CustomerID,
		AuditorUserID: uid,
		// 可选：把 disbursement_id、payment_channel_id 等也记录在 audits 里（更审计友好）
	}

	if err := auditDao.Create(ctx, record); err != nil {
		_ = tx.Rollback().Error
		response.Error(c, ecode.ErrCreateLoanAudits)
		return
	}

	// 9) commit
	if err := tx.Commit().Error; err != nil {
		logger.Error("tx commit failed", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InternalServerError)
		return
	}

	// 10) commit 成功后再 touch last_used_at
	go func(deviceID uint64) {
		// 创建独立上下文：脱离 gin 请求的 ctx，设置超时（防止无限阻塞）
		asyncCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel() // 确保超时/执行完后释放资源

		// 使用独立的 asyncCtx 执行数据库操作
		if err := h.userDao.TouchMFADeviceLastUsedAt(asyncCtx, deviceID); err != nil {
			logger.Warn("更新MFA设备最后使用时间失败", logger.Err(err), logger.Uint64("device_id", deviceID))
		}
	}(dev.ID) // 不再传递 gin 的 ctx

	response.Success(c, gin.H{})
}

// Create a new loanBaseinfo
// @Summary Create a new loanBaseinfo
// @Description Creates a new loanBaseinfo entity using the provided data in the request body.
// @Tags loanBaseinfo
// @Accept json
// @Produce json
// @Param data body types.CreateLoanBaseinfoRequest true "loanBaseinfo information"
// @Success 200 {object} types.CreateLoanBaseinfoReply{}
// @Router /api/v1/loanBaseinfo [post]
// @Security BearerAuth
func (h *loanBaseinfoHandler) Create(c *gin.Context) {
	form := &types.CreateLoanBaseinfoRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanBaseinfo := &model.LoanBaseinfo{}
	err = copier.Copy(loanBaseinfo, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanBaseinfo)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanBaseinfo)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanBaseinfo.ID})
}

// DeleteByID delete a loanBaseinfo by id
// @Summary Delete a loanBaseinfo by id
// @Description Deletes a existing loanBaseinfo identified by the given id in the path.
// @Tags loanBaseinfo
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanBaseinfoByIDReply{}
// @Router /api/v1/loanBaseinfo/{id} [delete]
// @Security BearerAuth
func (h *loanBaseinfoHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanBaseinfoIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByID(ctx, id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update a loanBaseinfo by id
// @Summary Update a loanBaseinfo by id
// @Description Updates the specified loanBaseinfo by given id in the path, support partial update.
// @Tags loanBaseinfo
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanBaseinfoByIDRequest true "loanBaseinfo information"
// @Success 200 {object} types.UpdateLoanBaseinfoByIDReply{}
// @Router /api/v1/loanBaseinfo/{id} [put]
// @Security BearerAuth
func (h *loanBaseinfoHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanBaseinfoIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanBaseinfoByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanBaseinfo := &model.LoanBaseinfo{}
	err = copier.Copy(loanBaseinfo, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanBaseinfo)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanBaseinfo)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanBaseinfo by id
// @Summary Get a loanBaseinfo by id
// @Description Gets detailed information of a loanBaseinfo specified by the given id in the path.
// @Tags loanBaseinfo
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanBaseinfoByIDReply{}
// @Router /api/v1/loanBaseinfo/{id} [get]
// @Security BearerAuth
func (h *loanBaseinfoHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanBaseinfoIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)

	loanBaseinfo, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.LoanBaseinfoObjDetail{}
	if err := copier.Copy(data, loanBaseinfo); err != nil {
		response.Error(c, ecode.ErrGetByIDLoanBaseinfo)
		return
	}

	// ✅ 查 files（按 type 分组）
	files, err := h.iDao.GetFilesMapByBaseinfoID(ctx, id)
	if err != nil {
		logger.Error("GetFilesMapByBaseinfoID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	// ✅ 保证返回的是 {} 而不是 null
	if files == nil {
		files = map[string][]string{}
	}
	data.Files = files

	response.Success(c, gin.H{"loanBaseinfo": data})
}

// List get a paginated list of loanBaseinfos by custom conditions
// @Summary Get a paginated list of loanBaseinfos by custom conditions
// @Description Returns a paginated list of loanBaseinfo based on query filters, including page number and size.
// @Tags loanBaseinfo
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanBaseinfosReply{}
// @Router /api/v1/loanBaseinfo/list [post]
// @Security BearerAuth
func (h *loanBaseinfoHandler) List(c *gin.Context) {
	form := &types.ListLoanBaseinfosRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanBaseinfos, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertSimpleLoanBaseinfos(loanBaseinfos)
	if err != nil {
		response.Error(c, ecode.ErrListLoanBaseinfo)
		return
	}

	response.Success(c, gin.H{
		"records": data,
		"total":   total,
	})
}

// DeleteByIDs batch delete loanBaseinfo by ids
// @Summary Batch delete loanBaseinfo by ids
// @Description Deletes multiple loanBaseinfo by a list of id
// @Tags loanBaseinfo
// @Param data body types.DeleteLoanBaseinfosByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanBaseinfosByIDsReply{}
// @Router /api/v1/loanBaseinfo/delete/ids [post]
// @Security BearerAuth
func (h *loanBaseinfoHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanBaseinfosByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err = h.iDao.DeleteByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByCondition get a loanBaseinfo by custom condition
// @Summary Get a loanBaseinfo by custom condition
// @Description Returns a single loanBaseinfo that matches the specified filter conditions.
// @Tags loanBaseinfo
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanBaseinfoByConditionReply{}
// @Router /api/v1/loanBaseinfo/condition [post]
// @Security BearerAuth
func (h *loanBaseinfoHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanBaseinfoByConditionRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	err = form.Conditions.CheckValid()
	if err != nil {
		logger.Warn("Parameters error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanBaseinfo, err := h.iDao.GetByCondition(ctx, &form.Conditions)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByCondition not found", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByCondition error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.LoanBaseinfoObjDetail{}
	err = copier.Copy(data, loanBaseinfo)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanBaseinfo)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanBaseinfo": data})
}

// ListByIDs batch get loanBaseinfo by ids
// @Summary Batch get loanBaseinfo by ids
// @Description Returns a list of loanBaseinfo that match the list of id.
// @Tags loanBaseinfo
// @Param data body types.ListLoanBaseinfosByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanBaseinfosByIDsReply{}
// @Router /api/v1/loanBaseinfo/list/ids [post]
// @Security BearerAuth
func (h *loanBaseinfoHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanBaseinfosByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanBaseinfoMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanBaseinfos := []*types.LoanBaseinfoObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanBaseinfoMap[id]; ok {
			record, err := convertLoanBaseinfo(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanBaseinfo)
				return
			}
			loanBaseinfos = append(loanBaseinfos, record)
		}
	}

	response.Success(c, gin.H{
		"loanBaseinfos": loanBaseinfos,
	})
}

// ListByLastID get a paginated list of loanBaseinfos by last id
// @Summary Get a paginated list of loanBaseinfos by last id
// @Description Returns a paginated list of loanBaseinfos starting after a given last id, useful for cursor-based pagination.
// @Tags loanBaseinfo
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanBaseinfosReply{}
// @Router /api/v1/loanBaseinfo/list [get]
// @Security BearerAuth
func (h *loanBaseinfoHandler) ListByLastID(c *gin.Context) {
	lastID := utils.StrToUint64(c.Query("lastID"))
	if lastID == 0 {
		lastID = math.MaxInt32
	}
	limit := utils.StrToInt(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	sort := c.Query("sort")

	ctx := middleware.WrapCtx(c)
	loanBaseinfos, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanBaseinfos(loanBaseinfos)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanBaseinfo)
		return
	}

	response.Success(c, gin.H{
		"loanBaseinfos": data,
	})
}

func getLoanBaseinfoIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertSimpleLoanBaseinfo(loanBaseinfo *model.LoanBaseinfo) (*types.LoanBaseinfoSimpleObjDetail, error) {
	data := &types.LoanBaseinfoSimpleObjDetail{}
	err := copier.Copy(data, loanBaseinfo)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertSimpleLoanBaseinfoWithAuditRecord(loanBaseinfo *model.LoanBaseinfoWithAuditRecord) (*types.LoanBaseinfoWithAuditRecords, error) {
	data := &types.LoanBaseinfoWithAuditRecords{}
	err := copier.Copy(data, loanBaseinfo)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanBaseinfo(loanBaseinfo *model.LoanBaseinfo) (*types.LoanBaseinfoObjDetail, error) {
	data := &types.LoanBaseinfoObjDetail{}
	err := copier.Copy(data, loanBaseinfo)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanBaseinfos(fromValues []*model.LoanBaseinfo) ([]*types.LoanBaseinfoObjDetail, error) {
	toValues := []*types.LoanBaseinfoObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanBaseinfo(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}
	return toValues, nil
}

func convertSimpleLoanBaseinfos(fromValues []*model.LoanBaseinfo) ([]*types.LoanBaseinfoSimpleObjDetail, error) {
	toValues := []*types.LoanBaseinfoSimpleObjDetail{}
	for _, v := range fromValues {
		data, err := convertSimpleLoanBaseinfo(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}
	return toValues, nil
}

func convertSimpleLoanBaseinfosWithAuditRecord(fromValues []*model.LoanBaseinfoWithAuditRecord) ([]*types.LoanBaseinfoWithAuditRecords, error) {
	toValues := []*types.LoanBaseinfoWithAuditRecords{}
	for _, v := range fromValues {
		data, err := convertSimpleLoanBaseinfoWithAuditRecord(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}
	return toValues, nil
}

func getCurrentTime() *time.Time {
	now := time.Now()
	return &now
}

func generateOrderNo() string {
	// 1. 获取当前时间，格式化为 年月日时分秒（例如：20260212153045）
	now := time.Now()
	timeStr := now.Format("20060102150405") // Go的时间格式化是固定参考时间：2006-01-02 15:04:05

	// 2. 生成3位随机数（避免同一秒内生成多个订单号导致重复）
	// 这里用纳秒取模，简单且无需额外依赖，适合单机场景
	randomNum := now.Nanosecond() % 1000 // 取纳秒的后3位，范围 0-999

	// 3. 拼接订单号：PO + 时间字符串 + 补零后的3位随机数
	orderNo := fmt.Sprintf("PO%s%03d", timeStr, randomNum)

	return orderNo
}
