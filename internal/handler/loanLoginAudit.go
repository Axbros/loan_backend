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

var _ LoanLoginAuditHandler = (*loanLoginAuditHandler)(nil)

// LoanLoginAuditHandler defining the handler interface
type LoanLoginAuditHandler interface {
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

type loanLoginAuditHandler struct {
	iDao dao.LoanLoginAuditDao
}

// NewLoanLoginAuditHandler creating the handler interface
func NewLoanLoginAuditHandler() LoanLoginAuditHandler {
	return &loanLoginAuditHandler{
		iDao: dao.NewLoanLoginAuditDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanLoginAuditCache(database.GetCacheType()),
		),
	}
}

// Create a new loanLoginAudit
// @Summary Create a new loanLoginAudit
// @Description Creates a new loanLoginAudit entity using the provided data in the request body.
// @Tags loanLoginAudit
// @Accept json
// @Produce json
// @Param data body types.CreateLoanLoginAuditRequest true "loanLoginAudit information"
// @Success 200 {object} types.CreateLoanLoginAuditReply{}
// @Router /api/v1/loanLoginAudit [post]
// @Security BearerAuth
func (h *loanLoginAuditHandler) Create(c *gin.Context) {
	form := &types.CreateLoanLoginAuditRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanLoginAudit := &model.LoanLoginAudit{}
	err = copier.Copy(loanLoginAudit, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanLoginAudit)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanLoginAudit)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanLoginAudit.ID})
}

// DeleteByID delete a loanLoginAudit by id
// @Summary Delete a loanLoginAudit by id
// @Description Deletes a existing loanLoginAudit identified by the given id in the path.
// @Tags loanLoginAudit
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanLoginAuditByIDReply{}
// @Router /api/v1/loanLoginAudit/{id} [delete]
// @Security BearerAuth
func (h *loanLoginAuditHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanLoginAuditIDFromPath(c)
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

// UpdateByID update a loanLoginAudit by id
// @Summary Update a loanLoginAudit by id
// @Description Updates the specified loanLoginAudit by given id in the path, support partial update.
// @Tags loanLoginAudit
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanLoginAuditByIDRequest true "loanLoginAudit information"
// @Success 200 {object} types.UpdateLoanLoginAuditByIDReply{}
// @Router /api/v1/loanLoginAudit/{id} [put]
// @Security BearerAuth
func (h *loanLoginAuditHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanLoginAuditIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanLoginAuditByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanLoginAudit := &model.LoanLoginAudit{}
	err = copier.Copy(loanLoginAudit, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanLoginAudit)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanLoginAudit)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanLoginAudit by id
// @Summary Get a loanLoginAudit by id
// @Description Gets detailed information of a loanLoginAudit specified by the given id in the path.
// @Tags loanLoginAudit
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanLoginAuditByIDReply{}
// @Router /api/v1/loanLoginAudit/{id} [get]
// @Security BearerAuth
func (h *loanLoginAuditHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanLoginAuditIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanLoginAudit, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanLoginAuditObjDetail{}
	err = copier.Copy(data, loanLoginAudit)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanLoginAudit)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanLoginAudit": data})
}

// List get a paginated list of loanLoginAudits by custom conditions
// @Summary Get a paginated list of loanLoginAudits by custom conditions
// @Description Returns a paginated list of loanLoginAudit based on query filters, including page number and size.
// @Tags loanLoginAudit
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanLoginAuditsReply{}
// @Router /api/v1/loanLoginAudit/list [post]
// @Security BearerAuth
func (h *loanLoginAuditHandler) List(c *gin.Context) {
	form := &types.ListLoanLoginAuditsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanLoginAudits, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanLoginAudits(loanLoginAudits)
	if err != nil {
		response.Error(c, ecode.ErrListLoanLoginAudit)
		return
	}

	response.Success(c, gin.H{
		"loanLoginAudits": data,
		"total":           total,
	})
}

// DeleteByIDs batch delete loanLoginAudit by ids
// @Summary Batch delete loanLoginAudit by ids
// @Description Deletes multiple loanLoginAudit by a list of id
// @Tags loanLoginAudit
// @Param data body types.DeleteLoanLoginAuditsByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanLoginAuditsByIDsReply{}
// @Router /api/v1/loanLoginAudit/delete/ids [post]
// @Security BearerAuth
func (h *loanLoginAuditHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanLoginAuditsByIDsRequest{}
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

// GetByCondition get a loanLoginAudit by custom condition
// @Summary Get a loanLoginAudit by custom condition
// @Description Returns a single loanLoginAudit that matches the specified filter conditions.
// @Tags loanLoginAudit
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanLoginAuditByConditionReply{}
// @Router /api/v1/loanLoginAudit/condition [post]
// @Security BearerAuth
func (h *loanLoginAuditHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanLoginAuditByConditionRequest{}
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
	loanLoginAudit, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanLoginAuditObjDetail{}
	err = copier.Copy(data, loanLoginAudit)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanLoginAudit)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanLoginAudit": data})
}

// ListByIDs batch get loanLoginAudit by ids
// @Summary Batch get loanLoginAudit by ids
// @Description Returns a list of loanLoginAudit that match the list of id.
// @Tags loanLoginAudit
// @Param data body types.ListLoanLoginAuditsByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanLoginAuditsByIDsReply{}
// @Router /api/v1/loanLoginAudit/list/ids [post]
// @Security BearerAuth
func (h *loanLoginAuditHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanLoginAuditsByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanLoginAuditMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanLoginAudits := []*types.LoanLoginAuditObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanLoginAuditMap[id]; ok {
			record, err := convertLoanLoginAudit(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanLoginAudit)
				return
			}
			loanLoginAudits = append(loanLoginAudits, record)
		}
	}

	response.Success(c, gin.H{
		"loanLoginAudits": loanLoginAudits,
	})
}

// ListByLastID get a paginated list of loanLoginAudits by last id
// @Summary Get a paginated list of loanLoginAudits by last id
// @Description Returns a paginated list of loanLoginAudits starting after a given last id, useful for cursor-based pagination.
// @Tags loanLoginAudit
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanLoginAuditsReply{}
// @Router /api/v1/loanLoginAudit/list [get]
// @Security BearerAuth
func (h *loanLoginAuditHandler) ListByLastID(c *gin.Context) {
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
	loanLoginAudits, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanLoginAudits(loanLoginAudits)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanLoginAudit)
		return
	}

	response.Success(c, gin.H{
		"loanLoginAudits": data,
	})
}

func getLoanLoginAuditIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanLoginAudit(loanLoginAudit *model.LoanLoginAudit) (*types.LoanLoginAuditObjDetail, error) {
	data := &types.LoanLoginAuditObjDetail{}
	err := copier.Copy(data, loanLoginAudit)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanLoginAudits(fromValues []*model.LoanLoginAudit) ([]*types.LoanLoginAuditObjDetail, error) {
	toValues := []*types.LoanLoginAuditObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanLoginAudit(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
