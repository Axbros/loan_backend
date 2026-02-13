package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"strconv"
	"strings"
	"time"

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
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	// 解析并校验审核类型
	auditType := parseAuditType(form.AuditType)
	if auditType == -1 {
		logger.Warn("无效的审核类型", logger.String("audit_type", strconv.Itoa(form.AuditType)), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidAuditType)
		return
	}

	// 查询当前用户是否启用 MFA
	u, err := h.userDao.GetByID(ctx, uid)
	if err != nil {
		logger.Error("查询用户信息失败", logger.Err(err), logger.Uint64("uid", uid), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGetByIDLoanUsers)
		return
	}

	// 核心审核逻辑函数
	// 核心审核逻辑函数
	auditCoreLogic := func() error {
		// 1. 查询放款基础信息
		loanBaseinfoRecord, err := h.iDao.GetByID(ctx, form.CustomerID)
		if err != nil {
			logger.Warn("GetByID LoanBaseinfo error: ", logger.Err(err), logger.Uint64("customer_id", form.CustomerID), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.ErrGetByIDLoanBaseinfo)
			return err
		}
		if loanBaseinfoRecord == nil {
			logger.Warn("放款基础信息不存在", logger.Uint64("customer_id", form.CustomerID), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.ErrGetByIDLoanBaseinfo)
			return errors.New("loan baseinfo not found")
		}

		// 2. 财务审核专属逻辑：仅审核通过时创建放款记录 + 还款计划
		if form.AuditType == FinanceReviewType {
			// 第一步：先判断审核结果，拒绝则直接跳过
			if !form.AuditResult {
				logger.Warn("财务审核拒绝放款，跳过创建放款记录", logger.Uint64("customer_id", form.CustomerID), middleware.GCtxRequestIDField(c))
			} else {
				if form.PaymentChannelID == 0 {
					response.Error(c, ecode.ErrPaymentChannel)
					return errors.New("payment channel id is 0")
				}

				// 查询支付渠道信息
				paymentChannelRecord, err := h.channelDao.GetByID(ctx, form.PaymentChannelID)
				if err != nil {
					logger.Error("查询支付渠道失败", logger.Err(err), logger.Uint64("channel_id", form.PaymentChannelID), middleware.GCtxRequestIDField(c))
					response.Error(c, ecode.ErrGetByIDLoanPaymentChannels)
					return err
				}
				if paymentChannelRecord == nil {
					logger.Warn("支付渠道不存在", logger.Uint64("channel_id", form.PaymentChannelID), middleware.GCtxRequestIDField(c))
					response.Error(c, ecode.ErrGetByIDLoanPaymentChannels)
					return errors.New("payment channel not found")
				}

				// 处理费率：空值默认0
				var feeRate decimal.Decimal
				if paymentChannelRecord.PayoutFeeRate != nil {
					feeRate = *paymentChannelRecord.PayoutFeeRate
				} else {
					feeRate = decimal.Zero
				}

				// 金额类型适配（decimal 类型计算）
				applicationAmount := *loanBaseinfoRecord.ApplicationAmount

				// 计算净金额：申请金额 * (1 - 费率)
				one := decimal.NewFromInt(1)
				discountRate := one.Sub(feeRate)
				netAmount := applicationAmount.Mul(discountRate)

				// 构造放款金额（*decimal.Decimal 类型）
				disburseAmount := applicationAmount // 放款金额 = 申请金额

				// 推荐人ID赋值（外键约束已删除，直接赋值）
				sourceReferrerUserID := loanBaseinfoRecord.ReferrerUserID

				// 获取当前时间
				now := time.Now()
				currentTime := &now

				// 构造放款记录
				disbursmentRecord := &model.LoanDisbursements{
					BaseinfoID:           form.CustomerID,
					DisburseAmount:       &disburseAmount,
					NetAmount:            &netAmount,
					Status:               1, // 已放款
					SourceReferrerUserID: sourceReferrerUserID,
					AuditorUserID:        uid,
					PayoutChannelID:      form.PaymentChannelID,
					AuditedAt:            currentTime,
					DisbursedAt:          currentTime,
					PayoutOrderNo:        generateOrderNo(),
				}

				// 创建放款记录（不使用事务，按你要求先完成逻辑）
				err = h.disbursmentDao.Create(ctx, disbursmentRecord)
				if err != nil {
					logger.Error("创建放款记录失败", logger.Err(err), logger.Uint64("customer_id", form.CustomerID), middleware.GCtxRequestIDField(c))
					response.Error(c, ecode.ErrCreateLoanDisbursements)
					return err
				}

				// ########### 生成还款计划（核心调整）###########
				// 1. 从loanBaseinfoRecord获取借款天数
				loanDays := loanBaseinfoRecord.LoanDays
				if loanDays <= 0 {
					logger.Warn("借款天数无效，默认设为30天", logger.Int("loan_days", loanDays), logger.Uint64("customer_id", form.CustomerID))
					loanDays = 30 // 兜底默认值
				}

				// 2. 计算应还日期：当前日期 + loanDays 天
				dueDate := now.AddDate(0, 0, loanDays)

				// ✅ 核心修正：正确创建仅含日期的 time.Time 实例
				// time.Date(年, 月, 日, 时, 分, 秒, 纳秒, 时区) → 时分秒设为0，只保留日期
				dueDateOnly := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, dueDate.Location())
				dueDatePtr := &dueDateOnly // 转为指针，适配模型的 *time.Time 类型

				// 3. 转换放款金额为分（元→分）
				disburseAmountCent := int(disburseAmount.Mul(decimal.NewFromInt(100)).IntPart())

				// 4. 构造还款计划记录
				scheduleRecord := &model.LoanRepaymentSchedules{
					DisbursementID: int64(disbursmentRecord.ID),
					InstallmentNo:  1,                  // 固定1期
					DueDate:        dueDatePtr,         // 当前日期+loanDays
					PrincipalDue:   disburseAmountCent, // 应还本金=放款金额
					InterestDue:    0,                  // 利息暂为0
					FeeDue:         disburseAmountCent, // 应还费用=放款金额
					PenaltyDue:     0,                  // 罚息暂为0
					TotalDue:       disburseAmountCent, // 应还总额=放款金额
					PaidPrincipal:  0,                  // 已还本金=0
					PaidInterest:   0,                  // 已还利息=0
					PaidFee:        0,                  // 已还费用=0
					PaidPenalty:    0,                  // 已还罚息=0
					PaidTotal:      0,                  // 已还总额=0
					Status:         0,                  // 未还清
					LastPaidAt:     nil,                // 最近还款时间=null
					SettledAt:      nil,                // 结清时间=null
				}

				// 5. 插入还款计划
				err = h.repaymentScheduleDao.Create(ctx, scheduleRecord)
				if err != nil {
					logger.Error("创建还款计划失败", logger.Err(err), logger.Uint64("disbursement_id", disbursmentRecord.ID), middleware.GCtxRequestIDField(c))
					response.Error(c, ecode.ErrCreateLoanRepaymentSchedules)
					return err
				}
			}
		}

		// 3. 更新放款基础信息审核状态（无论通过/拒绝都执行）
		var auditResult int
		if form.AuditResult {
			auditResult = 1
			loanBaseinfoRecord.AuditStatus = form.AuditType
		} else {
			auditResult = -1
			loanBaseinfoRecord.AuditStatus = -1
		}

		err = h.iDao.UpdateByID(ctx, loanBaseinfoRecord)
		if err != nil {
			logger.Warn("UpdateByID LoanBaseinfo error: ", logger.Err(err), logger.Uint64("customer_id", form.CustomerID), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.ErrUpdateByIDLoanBaseinfo)
			return err
		}

		// 4. 创建审核记录（无论通过/拒绝都执行）
		record := &model.LoanAudits{
			AuditResult:   auditResult,
			AuditType:     int(auditType),
			AuditComment:  form.Remark,
			BaseinfoID:    form.CustomerID,
			AuditorUserID: uid,
		}

		err = h.auditDao.Create(ctx, record)
		if err != nil {
			logger.Error("创建审核记录失败", logger.Err(err), logger.Uint64("customer_id", form.CustomerID), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.ErrCreateLoanAudits)
			return err
		}

		// 5. 返回成功
		response.Success(c, gin.H{})
		return nil
	}

	// MFA 启用场景
	if u.MfaEnabled == 1 {
		otpCode := strings.TrimSpace(form.MfaCode)
		if len(otpCode) != 6 {
			logger.Warn("MFA验证码长度错误", logger.String("otp_code", otpCode), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.MFAOTPRequired)
			return
		}

		// 查询用户的 active 主设备
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

		// 解密 secret
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

		// 校验 OTP
		ok, err := totp.ValidateCustom(otpCode, secret, time.Now(), totp.ValidateOpts{
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

		// 异步更新 last_used_at
		go func(ctx context.Context, deviceID uint64) {
			err := h.userDao.TouchMFADeviceLastUsedAt(ctx, deviceID)
			if err != nil {
				logger.Warn("更新MFA设备最后使用时间失败", logger.Err(err), logger.Uint64("device_id", deviceID))
			}
		}(ctx, dev.ID)

		// 执行核心审核逻辑
		if err := auditCoreLogic(); err != nil {
			return
		}
	} else {
		// MFA未启用时拒绝审核 ✅ 已按你的要求修改
		logger.Warn("用户未启用MFA，禁止审核", logger.Uint64("uid", uid), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.MFANotEnabled)
		return
	}
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
