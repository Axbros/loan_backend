package handler

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
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
	iDao     dao.LoanBaseinfoDao
	auditDao dao.LoanAuditsDao
	userDao  dao.LoanUsersDao
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
	PreReviewType     AuditType = iota //初审审核
	FinanceReviewType                  //放款审核
	IncomeReviewType                   //回款审核
)

func (h *loanBaseinfoHandler) WithAuditRecordList(c *gin.Context) {
	form := &types.ListLoanBaseinfosRequestWithAuditType{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)

	baseColumn := query.Column{
		Name:  "audit_status", // 首字母大写，匹配 Column 结构体定义
		Exp:   "=",            // 显式指定表达式为等于
		Value: PreReview,
		Logic: "&", // 显式指定逻辑为 AND
	}
	form.Columns = append(form.Columns, baseColumn)

	loanBaseinfos, total, err := h.iDao.GetByColumnsWithAuditRecords(ctx, &form.Params, form.AuditType)
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
	case 0:
		return PreReviewType
	case 1:
		return FinanceReviewType
	case 2:
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

	// ########### 新增：解析并校验请求的审核类型 ###########
	auditType := parseAuditType(form.AuditType)
	if auditType == 0 {
		//logger.Warn("parseAuditType error: ", logger.Err(err), "audit_type:", form.AuditType, middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InternalServerError) // 返无效审核类型错误
		return
	}
	// #######################################################

	// 0) 查当前用户是否启用 MFA
	u, err := h.userDao.GetByID(ctx, uid)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanUsers)
		return
	}
	if u.MfaEnabled == 1 {
		otpCode := strings.TrimSpace(form.MfaCode)
		if len(otpCode) != 6 {
			response.Error(c, ecode.MFAOTPRequired)
			return
		}

		// 1) 查用户的 active 主设备
		dev, err := h.userDao.GetActivePrimaryMFADevice(ctx, uid)
		if err != nil {
			response.Error(c, ecode.InternalServerError)
			return
		}

		// 2) 解密 secret
		secret, err := decryptSecretFromBytes(dev.SecretEnc)
		if err != nil {
			response.Error(c, ecode.ErrSecret)
			return
		}
		secret = strings.TrimSpace(secret)

		// 3) 校验 OTP（允许 1 个窗口偏移）
		ok, err := totp.ValidateCustom(otpCode, secret, time.Now(), totp.ValidateOpts{
			Period:    30,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		})
		if err != nil {
			response.Error(c, ecode.ErrSecret)
			return
		}
		if !ok {
			response.Error(c, ecode.InvalidOTP)
			return
		}

		// 4) 可选：更新 last_used_at
		_ = h.userDao.TouchMFADeviceLastUsedAt(ctx, dev.ID)

		// 先audit表插入数据
		var auditResult int

		loanBaseinfoRecord, err := h.iDao.GetByID(ctx, form.CustomerID)
		if err != nil {
			logger.Warn("GetByID error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.ErrGetByIDLoanBaseinfo)
			return
		}
		// 审核通过/驳回
		if form.AuditResult {
			auditResult = 1
		} else {
			auditResult = -1
		}
		// 简化赋值：无需两次写，一次赋值即可
		loanBaseinfoRecord.AuditStatus = auditResult

		err = h.iDao.UpdateByID(ctx, loanBaseinfoRecord)
		if err != nil {
			logger.Warn("UpdateByID error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.ErrUpdateByIDLoanBaseinfo)
			return
		}

		// ########### 核心修复：替换写死的PreReviewType为解析后的auditType ###########
		record := &model.LoanAudits{
			AuditResult:   auditResult,
			AuditType:     int(auditType), // 用请求传的审核类型
			AuditComment:  form.Remark,
			BaseinfoID:    form.CustomerID,
			AuditorUserID: uid,
		}
		// ###########################################################################

		err = h.auditDao.Create(ctx, record)
		if err != nil {
			response.Error(c, ecode.ErrCreateLoanAudits)
			return
		}
		response.Success(c, gin.H{})
	} else {
		// ########### 新增：补充MFA未启用的分支处理（原代码缺失，会导致无响应） ###########
		logger.Warn("user MFA not enabled, uid:"+utils.Uint64ToStr(uid), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.MFANotEnabled) // 需在ecode中定义：MFA未启用
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
