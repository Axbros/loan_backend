package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"loan/internal/config"
	"loan/internal/tool"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"
	"github.com/google/uuid"

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
	PreReview(c *gin.Context)
	FinanceReview(c *gin.Context)
	WithAuditRecordList(c *gin.Context)
	UploadCertificate(c *gin.Context)
	GetCertificateBase64(c *gin.Context)
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

type AuditType int // 修正原Audit_Type命名，符合Go大驼峰规范
const (
	PreReviewType     = 1 //初审审核
	FinanceReviewType = 2 //放款审核
	//IncomeReviewType  = 3 //回款审核
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

func (h *loanBaseinfoHandler) PreReview(c *gin.Context) {
	h.review(c, PreReviewType)
}

func (h *loanBaseinfoHandler) FinanceReview(c *gin.Context) {
	h.review(c, FinanceReviewType)
}

func (h *loanBaseinfoHandler) review(c *gin.Context, forcedAuditType int) {
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

	// 不再使用前端传入的 AuditType，统一由路由/handler 决定
	auditType := forcedAuditType

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
	ok, err = tool.ValidateMFA(c, uid, otpCode)
	if err != nil || !ok {
		logger.Warn("ValidateMFA error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return
	}

	// 2) 开事务
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

	// 3) 锁住 baseinfo 行
	loanBaseinfoRecord := &model.LoanBaseinfo{}
	if err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", form.CustomerID).
		First(loanBaseinfoRecord).Error; err != nil {
		_ = tx.Rollback().Error
		logger.Warn("baseinfo not found / lock failed", logger.Err(err), logger.Uint64("customer_id", form.CustomerID), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGetByIDLoanBaseinfo)
		return
	}

	// 4) 状态机校验
	if loanBaseinfoRecord.AuditStatus == -1 {
		_ = tx.Rollback().Error
		response.Error(c, ecode.InvalidParams)
		return
	}

	// 建议加顺序校验
	if auditType == FinanceReviewType && loanBaseinfoRecord.AuditStatus != PreReviewType {
		_ = tx.Rollback().Error
		response.Error(c, ecode.InvalidParams)
		return
	}

	// 如果不允许重复初审通过，也可加
	// if auditType == PreReviewType && loanBaseinfoRecord.AuditStatus == PreReviewType {
	//     _ = tx.Rollback().Error
	//     response.Error(c, ecode.InvalidParams)
	//     return
	// }

	// 5) 财务审核通过：放款 + 还款计划
	var createdDisbursementID uint64
	if auditType == FinanceReviewType && form.AuditResult {
		if form.PaymentChannelID == 0 {
			_ = tx.Rollback().Error
			response.Error(c, ecode.ErrPaymentChannel)
			return
		}

		paymentChannelRecord, err := h.channelDao.GetByID(ctx, form.PaymentChannelID)
		if err != nil || paymentChannelRecord == nil {
			_ = tx.Rollback().Error
			response.Error(c, ecode.ErrGetByIDLoanPaymentChannels)
			return
		}

		existing := &model.LoanDisbursements{}
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("baseinfo_id = ?", form.CustomerID).
			Order("id DESC").
			First(existing).Error

		if err == nil && existing != nil && existing.ID != 0 {
			createdDisbursementID = existing.ID
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			if loanBaseinfoRecord.ApplicationAmount == 0 {
				_ = tx.Rollback().Error
				response.Error(c, ecode.InvalidParams)
				return
			}

			var feeRate = 0
			if paymentChannelRecord.PayoutFeeRate != 0 {
				feeRate = paymentChannelRecord.PayoutFeeRate
			}

			applicationAmount := loanBaseinfoRecord.ApplicationAmount

			feeAmount := int64(float64(applicationAmount) * float64(feeRate) / 100)
			netAmount := applicationAmount - feeAmount
			disburseAmount := applicationAmount

			now := time.Now()
			currentTime := &now

			disbursmentRecord := &model.LoanDisbursements{
				BaseinfoID:           form.CustomerID,
				DisburseAmount:       disburseAmount,
				NetAmount:            netAmount,
				Status:               1,
				SourceReferrerUserID: loanBaseinfoRecord.ReferrerUserID,
				AuditorUserID:        uid,
				PayoutChannelID:      form.PaymentChannelID,
				AuditedAt:            currentTime,
				DisbursedAt:          currentTime,
				PayoutOrderNo:        generateOrderNo("PO"),
			}

			if _, err := h.disbursmentDao.CreateByTx(ctx, tx, disbursmentRecord); err != nil {
				_ = tx.Rollback().Error
				response.Error(c, ecode.ErrCreateLoanDisbursements)
				return
			}
			createdDisbursementID = disbursmentRecord.ID

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

			if _, err := h.repaymentScheduleDao.CreateByTx(ctx, tx, scheduleRecord); err != nil {
				_ = tx.Rollback().Error
				response.Error(c, ecode.ErrCreateLoanRepaymentSchedules)
				return
			}
		} else {
			_ = tx.Rollback().Error
			response.Error(c, ecode.InternalServerError)
			return
		}
	}

	// 6) 更新审核状态
	auditResult := -1
	if form.AuditResult {
		auditResult = 1
		loanBaseinfoRecord.AuditStatus = auditType
	} else {
		loanBaseinfoRecord.AuditStatus = -1
	}

	if err := h.iDao.UpdateByTx(ctx, tx, loanBaseinfoRecord); err != nil {
		_ = tx.Rollback().Error
		response.Error(c, ecode.ErrUpdateByIDLoanBaseinfo)
		return
	}

	// 7) 写审核记录
	record := &model.LoanAudits{
		AuditResult:   auditResult,
		AuditType:     auditType,
		AuditComment:  form.Remark,
		BaseinfoID:    form.CustomerID,
		AuditorUserID: uid,
	}

	if _, err := h.auditDao.CreateByTx(ctx, tx, record); err != nil {
		_ = tx.Rollback().Error
		response.Error(c, ecode.ErrCreateLoanAudits)
		return
	}

	// 8) commit
	if err := tx.Commit().Error; err != nil {
		logger.Error("tx commit failed", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InternalServerError)
		return
	}

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
	// If ReferrerUserID is 0, set it to nil to avoid foreign key constraint error
	if form.ReferrerUserID == 0 {
		loanBaseinfo.ReferrerUserID = nil
	} else {
		uid := form.ReferrerUserID
		loanBaseinfo.ReferrerUserID = &uid
	}

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

// UploadCertificate upload baseinfo certificate
func (h *loanBaseinfoHandler) UploadCertificate(c *gin.Context) {
	// 1. Get file from form (field name: certificate)
	file, fileHeader, err := c.Request.FormFile("certificate")
	if err != nil {
		response.Error(c, ecode.InvalidParams)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	// 2. Validate file type (png/jpg/jpeg)
	allowedExts := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
	}
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExts[fileExt] {
		response.Error(c, ecode.ErrUnsupportedFileTypeBaseinfo)
		return
	}

	// 3. Ensure storage directory exists
	storageDir := config.Get().Storage.BaseinfoCertificate
	cleanStorageDir := filepath.Clean(storageDir)

	if err := os.MkdirAll(cleanStorageDir, 0755); err != nil {
		response.Error(c, ecode.ErrCreateFileFolderBaseinfo)
		return
	}

	// 4. Generate unique filename
	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("%s%s", uniqueID, fileExt)
	filePath := filepath.Join(storageDir, filename)

	// 5. Create local file
	cleanFilePath := filepath.Clean(filePath)

	if !strings.HasPrefix(cleanFilePath, cleanStorageDir) {
		logger.Errorf(
			"[UploadCertificate] invalid file path: filePath=%q cleanFilePath=%q cleanStorageDir=%q",
			filePath,
			cleanFilePath,
			cleanStorageDir,
		)

		response.Error(c, ecode.ErrInvalidFilePathBaseinfo)
		return
	}
	dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		response.Error(c, ecode.ErrSaveFileBaseinfo)
		return
	}
	defer func() {
		if dstFile != nil {
			_ = dstFile.Close()
		}
	}()

	// 6. Copy file content
	_, err = io.Copy(dstFile, file)
	if err != nil {
		response.Error(c, ecode.ErrSaveFileBaseinfo)
		return
	}

	// 7. Return success result
	response.Success(c, gin.H{
		"file_name": filename,
		"size":      fileHeader.Size,
	})
}

// GetCertificateBase64 get baseinfo certificate base64
func (h *loanBaseinfoHandler) GetCertificateBase64(c *gin.Context) {
	// 1. Get file_name param
	fileName := c.Param("file_name")
	if fileName == "" {
		response.Error(c, ecode.InvalidParams)
		return
	}

	if fileName != filepath.Base(fileName) {
		response.Error(c, ecode.ErrInvalidFilePathBaseinfo)
		return
	}

	// 2. Define storage directory
	storageDir := config.Get().Storage.BaseinfoCertificate
	cleanStorageDir := filepath.Clean(storageDir)

	// 3. Join path and validate
	filePath := filepath.Join(storageDir, fileName)
	cleanFilePath := filepath.Clean(filePath)
	if !strings.HasPrefix(cleanFilePath, cleanStorageDir) {
		response.Error(c, ecode.ErrInvalidFilePathBaseinfo)
		return
	}

	// 4. Check if file exists
	if _, err := os.Stat(cleanFilePath); os.IsNotExist(err) {
		response.Error(c, ecode.ErrFileNotFoundBaseinfo)
		return
	} else if err != nil {
		response.Error(c, ecode.ErrReadFileBaseinfo)
		return
	}

	// 5. Read file content
	fileContent, err := os.ReadFile(cleanFilePath)
	if err != nil {
		response.Error(c, ecode.ErrReadFileBaseinfo)
		return
	}

	// 6. Encode to Base64
	base64Str := base64.StdEncoding.EncodeToString(fileContent)

	// 7. Generate Data URI
	fileExt := strings.ToLower(filepath.Ext(fileName))
	mimeType := ""
	switch fileExt {
	case ".png":
		mimeType = "image/png"
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	default:
		mimeType = "application/octet-stream"
	}
	base64WithPrefix := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)

	// 8. Return result
	response.Success(c, gin.H{
		"base64_with_prefix": base64WithPrefix,
		"file_size":          len(fileContent),
	})
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

func generateOrderNo(prefix string) string {
	// 1. 获取当前时间，格式化为 年月日时分秒（例如：20260212153045）
	now := time.Now()
	timeStr := now.Format("20060102150405") // Go的时间格式化是固定参考时间：2006-01-02 15:04:05

	// 2. 生成3位随机数（避免同一秒内生成多个订单号导致重复）
	// 这里用纳秒取模，简单且无需额外依赖，适合单机场景
	randomNum := now.Nanosecond() % 1000 // 取纳秒的后3位，范围 0-999

	// 3. 拼接订单号：PO + 时间字符串 + 补零后的3位随机数
	orderNo := fmt.Sprintf("%s%s%03d", prefix, timeStr, randomNum)

	return orderNo
}
