package handler

import (
	"errors"
	"math"
	"time"

	"github.com/gin-gonic/gin"

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

var _ LoanCollectionLogsHandler = (*loanCollectionLogsHandler)(nil)

// LoanCollectionLogsHandler defining the handler interface
type LoanCollectionLogsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)

	DeleteByIDs(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	ListByLastID(c *gin.Context)
}

type loanCollectionLogsHandler struct {
	iDao    dao.LoanCollectionLogsDao
	caseDao dao.LoanCollectionCasesDao
}

// NewLoanCollectionLogsHandler creating the handler interface
func NewLoanCollectionLogsHandler() LoanCollectionLogsHandler {
	return &loanCollectionLogsHandler{
		iDao: dao.NewLoanCollectionLogsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanCollectionLogsCache(database.GetCacheType()),
		),
		caseDao: dao.NewLoanCollectionCasesDao(
			database.GetDB(),
			cache.NewLoanCollectionCasesCache(database.GetCacheType()),
		),
	}
}

// Create a new loanCollectionLogs
// @Summary Create a new loanCollectionLogs
// @Description Creates a new loanCollectionLogs entity using the provided data in the request body.
// @Tags loanCollectionLogs
// @Accept json
// @Produce json
// @Param data body types.CreateLoanCollectionLogsRequest true "loanCollectionLogs information"
// @Success 200 {object} types.CreateLoanCollectionLogsReply{}
// @Router /api/v1/loanCollectionLogs [post]
// @Security BearerAuth
func (h *loanCollectionLogsHandler) Create(c *gin.Context) {
	uid, ok := getUIDFromClaims(c)
	if !ok || uid == 0 {
		response.Error(c, ecode.Unauthorized)
		return
	}

	form := &types.CreateLoanCollectionLogsRequest{}
	if err := c.ShouldBindJSON(form); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanCollectionLogs := &model.LoanCollectionLogs{}

	// copier 只拷贝“同类型/可直接赋值”的字段（ActionType/Content 等）
	if err := copier.Copy(loanCollectionLogs, form); err != nil {
		logger.Error("copier.Copy error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrCreateLoanCollectionLogs)
		return
	}

	// 强制字段：不信任前端
	loanCollectionLogs.CollectorUserID = uid

	// CaseID: int64 -> uint64，手动转（并做校验）
	if form.CaseID <= 0 {
		response.Error(c, ecode.InvalidParams)
		return
	}
	loanCollectionLogs.CaseID = form.CaseID

	// NextFollowUpAt: string -> *time.Time，手动 parse
	if form.NextFollowUpAt != "" {
		t, err := time.Parse("2006-01-02 15:04:05", form.NextFollowUpAt)
		if err != nil {
			logger.Warn("Parse NextFollowUpAt error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.InvalidParams)
			return
		}
		loanCollectionLogs.NextFollowUpAt = &t
	} else {
		loanCollectionLogs.NextFollowUpAt = nil
	}

	ctx := middleware.WrapCtx(c)

	record, err := h.caseDao.GetByID(ctx, loanCollectionLogs.CaseID)
	if err != nil {
		logger.Error("GetByID error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGetByIDLoanCollectionCases)
		return
	}
	if record.Status == 2 {
		//最终态禁止更新
		response.Error(c, ecode.InternalServerError)
		return
	}
	if err := h.iDao.Create(ctx, loanCollectionLogs); err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	record.Status = 1 //重要调用创建催收日志 就一定是 跟进中的状态

	err = h.caseDao.UpdateByID(ctx, record)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanCollectionLogs)
		return
	}

	response.Success(c, gin.H{"id": loanCollectionLogs.ID})
}

// DeleteByID delete a loanCollectionLogs by id
// @Summary Delete a loanCollectionLogs by id
// @Description Deletes a existing loanCollectionLogs identified by the given id in the path.
// @Tags loanCollectionLogs
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanCollectionLogsByIDReply{}
// @Router /api/v1/loanCollectionLogs/{id} [delete]
// @Security BearerAuth
func (h *loanCollectionLogsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanCollectionLogsIDFromPath(c)
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

// UpdateByID update a loanCollectionLogs by id
// @Summary Update a loanCollectionLogs by id
// @Description Updates the specified loanCollectionLogs by given id in the path, support partial update.
// @Tags loanCollectionLogs
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanCollectionLogsByIDRequest true "loanCollectionLogs information"
// @Success 200 {object} types.UpdateLoanCollectionLogsByIDReply{}
// @Router /api/v1/loanCollectionLogs/{id} [put]
// @Security BearerAuth
func (h *loanCollectionLogsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanCollectionLogsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanCollectionLogsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanCollectionLogs := &model.LoanCollectionLogs{}
	err = copier.Copy(loanCollectionLogs, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanCollectionLogs)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanCollectionLogs)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanCollectionLogs by id
// @Summary Get a loanCollectionLogs by id
// @Description Gets detailed information of a loanCollectionLogs specified by the given id in the path.
// @Tags loanCollectionLogs
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanCollectionLogsByIDReply{}
// @Router /api/v1/loanCollectionLogs/{id} [get]
// @Security BearerAuth
func (h *loanCollectionLogsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanCollectionLogsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanCollectionLogs, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanCollectionLogsObjDetail{}
	err = copier.Copy(data, loanCollectionLogs)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanCollectionLogs)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanCollectionLogs": data})
}

// List get a paginated list of loanCollectionLogss by custom conditions
// @Summary Get a paginated list of loanCollectionLogss by custom conditions
// @Description Returns a paginated list of loanCollectionLogs based on query filters, including page number and size.
// @Tags loanCollectionLogs
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanCollectionLogssReply{}
// @Router /api/v1/loanCollectionLogs/list [post]
// @Security BearerAuth
func (h *loanCollectionLogsHandler) List(c *gin.Context) {
	form := &types.ListLoanCollectionLogssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanCollectionLogss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanCollectionLogss(loanCollectionLogss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanCollectionLogs)
		return
	}

	response.Success(c, gin.H{
		"records": data,
		"total":   total,
	})
}

// DeleteByIDs batch delete loanCollectionLogs by ids
// @Summary Batch delete loanCollectionLogs by ids
// @Description Deletes multiple loanCollectionLogs by a list of id
// @Tags loanCollectionLogs
// @Param data body types.DeleteLoanCollectionLogssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanCollectionLogssByIDsReply{}
// @Router /api/v1/loanCollectionLogs/delete/ids [post]
// @Security BearerAuth
func (h *loanCollectionLogsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanCollectionLogssByIDsRequest{}
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

// GetByCondition get a loanCollectionLogs by custom condition
// @Summary Get a loanCollectionLogs by custom condition
// @Description Returns a single loanCollectionLogs that matches the specified filter conditions.
// @Tags loanCollectionLogs
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanCollectionLogsByConditionReply{}
// @Router /api/v1/loanCollectionLogs/condition [post]
// @Security BearerAuth
func (h *loanCollectionLogsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanCollectionLogsByConditionRequest{}
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
	loanCollectionLogs, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanCollectionLogsObjDetail{}
	err = copier.Copy(data, loanCollectionLogs)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanCollectionLogs)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanCollectionLogs": data})
}

// ListByIDs batch get loanCollectionLogs by ids
// @Summary Batch get loanCollectionLogs by ids
// @Description Returns a list of loanCollectionLogs that match the list of id.
// @Tags loanCollectionLogs
// @Param data body types.ListLoanCollectionLogssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanCollectionLogssByIDsReply{}
// @Router /api/v1/loanCollectionLogs/list/ids [post]
// @Security BearerAuth
func (h *loanCollectionLogsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanCollectionLogssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanCollectionLogsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanCollectionLogss := []*types.LoanCollectionLogsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanCollectionLogsMap[id]; ok {
			record, err := convertLoanCollectionLogs(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanCollectionLogs)
				return
			}
			loanCollectionLogss = append(loanCollectionLogss, record)
		}
	}

	response.Success(c, gin.H{
		"loanCollectionLogss": loanCollectionLogss,
	})
}

// ListByLastID get a paginated list of loanCollectionLogss by last id
// @Summary Get a paginated list of loanCollectionLogss by last id
// @Description Returns a paginated list of loanCollectionLogss starting after a given last id, useful for cursor-based pagination.
// @Tags loanCollectionLogs
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanCollectionLogssReply{}
// @Router /api/v1/loanCollectionLogs/list [get]
// @Security BearerAuth
func (h *loanCollectionLogsHandler) ListByLastID(c *gin.Context) {
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
	loanCollectionLogss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanCollectionLogss(loanCollectionLogss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanCollectionLogs)
		return
	}

	response.Success(c, gin.H{
		"loanCollectionLogss": data,
	})
}

func getLoanCollectionLogsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanCollectionLogs(loanCollectionLogs *model.LoanCollectionLogs) (*types.LoanCollectionLogsObjDetail, error) {
	data := &types.LoanCollectionLogsObjDetail{}
	err := copier.Copy(data, loanCollectionLogs)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanCollectionLogss(fromValues []*model.LoanCollectionLogs) ([]*types.LoanCollectionLogsObjDetail, error) {
	toValues := []*types.LoanCollectionLogsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanCollectionLogs(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
