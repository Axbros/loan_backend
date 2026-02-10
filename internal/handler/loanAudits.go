package handler

import (
	"errors"
	"math"

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

var _ LoanAuditsHandler = (*loanAuditsHandler)(nil)

// LoanAuditsHandler defining the handler interface
type LoanAuditsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Detail(c *gin.Context)

	DeleteByIDs(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	ListByLastID(c *gin.Context)
}

type loanAuditsHandler struct {
	iDao dao.LoanAuditsDao
}

// NewLoanAuditsHandler creating the handler interface
func NewLoanAuditsHandler() LoanAuditsHandler {
	return &loanAuditsHandler{
		iDao: dao.NewLoanAuditsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanAuditsCache(database.GetCacheType()),
		),
	}
}

func (h *loanAuditsHandler) Detail(c *gin.Context) {
	form := &types.RequestLoanAuditsDetail{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)
	record, err := h.iDao.GetByBaseinfoID(ctx, form.BaseinfoID, form.AuditType)
	if err != nil {
		logger.Warn("GetByBaseinfoID error: ", logger.Err(err))
		response.Error(c, ecode.InternalServerError)
		return
	}
	response.Success(c, gin.H{
		"record": record,
	})

}

// Create a new loanAudits
// @Summary Create a new loanAudits
// @Description Creates a new loanAudits entity using the provided data in the request body.
// @Tags loanAudits
// @Accept json
// @Produce json
// @Param data body types.CreateLoanAuditsRequest true "loanAudits information"
// @Success 200 {object} types.CreateLoanAuditsReply{}
// @Router /api/v1/loanAudits [post]
// @Security BearerAuth
func (h *loanAuditsHandler) Create(c *gin.Context) {
	form := &types.CreateLoanAuditsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanAudits := &model.LoanAudits{}
	err = copier.Copy(loanAudits, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanAudits)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanAudits)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanAudits.ID})
}

// DeleteByID delete a loanAudits by id
// @Summary Delete a loanAudits by id
// @Description Deletes a existing loanAudits identified by the given id in the path.
// @Tags loanAudits
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanAuditsByIDReply{}
// @Router /api/v1/loanAudits/{id} [delete]
// @Security BearerAuth
func (h *loanAuditsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanAuditsIDFromPath(c)
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

// UpdateByID update a loanAudits by id
// @Summary Update a loanAudits by id
// @Description Updates the specified loanAudits by given id in the path, support partial update.
// @Tags loanAudits
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanAuditsByIDRequest true "loanAudits information"
// @Success 200 {object} types.UpdateLoanAuditsByIDReply{}
// @Router /api/v1/loanAudits/{id} [put]
// @Security BearerAuth
func (h *loanAuditsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanAuditsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanAuditsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanAudits := &model.LoanAudits{}
	err = copier.Copy(loanAudits, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanAudits)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanAudits)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanAudits by id
// @Summary Get a loanAudits by id
// @Description Gets detailed information of a loanAudits specified by the given id in the path.
// @Tags loanAudits
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanAuditsByIDReply{}
// @Router /api/v1/loanAudits/{id} [get]
// @Security BearerAuth
func (h *loanAuditsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanAuditsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanAudits, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanAuditsObjDetail{}
	err = copier.Copy(data, loanAudits)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanAudits)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanAudits": data})
}

// List get a paginated list of loanAuditss by custom conditions
// @Summary Get a paginated list of loanAuditss by custom conditions
// @Description Returns a paginated list of loanAudits based on query filters, including page number and size.
// @Tags loanAudits
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanAuditssReply{}
// @Router /api/v1/loanAudits/list [post]
// @Security BearerAuth
func (h *loanAuditsHandler) List(c *gin.Context) {
	form := &types.ListLoanAuditssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanAuditss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanAuditss(loanAuditss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanAudits)
		return
	}

	response.Success(c, gin.H{
		"loanAuditss": data,
		"total":       total,
	})
}

// DeleteByIDs batch delete loanAudits by ids
// @Summary Batch delete loanAudits by ids
// @Description Deletes multiple loanAudits by a list of id
// @Tags loanAudits
// @Param data body types.DeleteLoanAuditssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanAuditssByIDsReply{}
// @Router /api/v1/loanAudits/delete/ids [post]
// @Security BearerAuth
func (h *loanAuditsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanAuditssByIDsRequest{}
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

// GetByCondition get a loanAudits by custom condition
// @Summary Get a loanAudits by custom condition
// @Description Returns a single loanAudits that matches the specified filter conditions.
// @Tags loanAudits
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanAuditsByConditionReply{}
// @Router /api/v1/loanAudits/condition [post]
// @Security BearerAuth
func (h *loanAuditsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanAuditsByConditionRequest{}
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
	loanAudits, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanAuditsObjDetail{}
	err = copier.Copy(data, loanAudits)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanAudits)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanAudits": data})
}

// ListByIDs batch get loanAudits by ids
// @Summary Batch get loanAudits by ids
// @Description Returns a list of loanAudits that match the list of id.
// @Tags loanAudits
// @Param data body types.ListLoanAuditssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanAuditssByIDsReply{}
// @Router /api/v1/loanAudits/list/ids [post]
// @Security BearerAuth
func (h *loanAuditsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanAuditssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanAuditsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanAuditss := []*types.LoanAuditsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanAuditsMap[id]; ok {
			record, err := convertLoanAudits(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanAudits)
				return
			}
			loanAuditss = append(loanAuditss, record)
		}
	}

	response.Success(c, gin.H{
		"loanAuditss": loanAuditss,
	})
}

// ListByLastID get a paginated list of loanAuditss by last id
// @Summary Get a paginated list of loanAuditss by last id
// @Description Returns a paginated list of loanAuditss starting after a given last id, useful for cursor-based pagination.
// @Tags loanAudits
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanAuditssReply{}
// @Router /api/v1/loanAudits/list [get]
// @Security BearerAuth
func (h *loanAuditsHandler) ListByLastID(c *gin.Context) {
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
	loanAuditss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanAuditss(loanAuditss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanAudits)
		return
	}

	response.Success(c, gin.H{
		"loanAuditss": data,
	})
}

func getLoanAuditsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanAudits(loanAudits *model.LoanAudits) (*types.LoanAuditsObjDetail, error) {
	data := &types.LoanAuditsObjDetail{}
	err := copier.Copy(data, loanAudits)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanAuditss(fromValues []*model.LoanAudits) ([]*types.LoanAuditsObjDetail, error) {
	toValues := []*types.LoanAuditsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanAudits(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
