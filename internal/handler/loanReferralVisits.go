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

var _ LoanReferralVisitsHandler = (*loanReferralVisitsHandler)(nil)

// LoanReferralVisitsHandler defining the handler interface
type LoanReferralVisitsHandler interface {
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

type loanReferralVisitsHandler struct {
	iDao dao.LoanReferralVisitsDao
}

// NewLoanReferralVisitsHandler creating the handler interface
func NewLoanReferralVisitsHandler() LoanReferralVisitsHandler {
	return &loanReferralVisitsHandler{
		iDao: dao.NewLoanReferralVisitsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanReferralVisitsCache(database.GetCacheType()),
		),
	}
}

// Create a new loanReferralVisits
// @Summary Create a new loanReferralVisits
// @Description Creates a new loanReferralVisits entity using the provided data in the request body.
// @Tags loanReferralVisits
// @Accept json
// @Produce json
// @Param data body types.CreateLoanReferralVisitsRequest true "loanReferralVisits information"
// @Success 200 {object} types.CreateLoanReferralVisitsReply{}
// @Router /api/v1/loanReferralVisits [post]
// @Security BearerAuth
func (h *loanReferralVisitsHandler) Create(c *gin.Context) {
	form := &types.CreateLoanReferralVisitsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanReferralVisits := &model.LoanReferralVisits{}
	err = copier.Copy(loanReferralVisits, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanReferralVisits)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanReferralVisits)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanReferralVisits.ID})
}

// DeleteByID delete a loanReferralVisits by id
// @Summary Delete a loanReferralVisits by id
// @Description Deletes a existing loanReferralVisits identified by the given id in the path.
// @Tags loanReferralVisits
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanReferralVisitsByIDReply{}
// @Router /api/v1/loanReferralVisits/{id} [delete]
// @Security BearerAuth
func (h *loanReferralVisitsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanReferralVisitsIDFromPath(c)
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

// UpdateByID update a loanReferralVisits by id
// @Summary Update a loanReferralVisits by id
// @Description Updates the specified loanReferralVisits by given id in the path, support partial update.
// @Tags loanReferralVisits
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanReferralVisitsByIDRequest true "loanReferralVisits information"
// @Success 200 {object} types.UpdateLoanReferralVisitsByIDReply{}
// @Router /api/v1/loanReferralVisits/{id} [put]
// @Security BearerAuth
func (h *loanReferralVisitsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanReferralVisitsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanReferralVisitsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanReferralVisits := &model.LoanReferralVisits{}
	err = copier.Copy(loanReferralVisits, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanReferralVisits)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanReferralVisits)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanReferralVisits by id
// @Summary Get a loanReferralVisits by id
// @Description Gets detailed information of a loanReferralVisits specified by the given id in the path.
// @Tags loanReferralVisits
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanReferralVisitsByIDReply{}
// @Router /api/v1/loanReferralVisits/{id} [get]
// @Security BearerAuth
func (h *loanReferralVisitsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanReferralVisitsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanReferralVisits, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanReferralVisitsObjDetail{}
	err = copier.Copy(data, loanReferralVisits)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanReferralVisits)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanReferralVisits": data})
}

// List get a paginated list of loanReferralVisitss by custom conditions
// @Summary Get a paginated list of loanReferralVisitss by custom conditions
// @Description Returns a paginated list of loanReferralVisits based on query filters, including page number and size.
// @Tags loanReferralVisits
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanReferralVisitssReply{}
// @Router /api/v1/loanReferralVisits/list [post]
// @Security BearerAuth
func (h *loanReferralVisitsHandler) List(c *gin.Context) {
	form := &types.ListLoanReferralVisitssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanReferralVisitss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanReferralVisitss(loanReferralVisitss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanReferralVisits)
		return
	}

	response.Success(c, gin.H{
		"loanReferralVisitss": data,
		"total":               total,
	})
}

// DeleteByIDs batch delete loanReferralVisits by ids
// @Summary Batch delete loanReferralVisits by ids
// @Description Deletes multiple loanReferralVisits by a list of id
// @Tags loanReferralVisits
// @Param data body types.DeleteLoanReferralVisitssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanReferralVisitssByIDsReply{}
// @Router /api/v1/loanReferralVisits/delete/ids [post]
// @Security BearerAuth
func (h *loanReferralVisitsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanReferralVisitssByIDsRequest{}
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

// GetByCondition get a loanReferralVisits by custom condition
// @Summary Get a loanReferralVisits by custom condition
// @Description Returns a single loanReferralVisits that matches the specified filter conditions.
// @Tags loanReferralVisits
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanReferralVisitsByConditionReply{}
// @Router /api/v1/loanReferralVisits/condition [post]
// @Security BearerAuth
func (h *loanReferralVisitsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanReferralVisitsByConditionRequest{}
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
	loanReferralVisits, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanReferralVisitsObjDetail{}
	err = copier.Copy(data, loanReferralVisits)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanReferralVisits)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanReferralVisits": data})
}

// ListByIDs batch get loanReferralVisits by ids
// @Summary Batch get loanReferralVisits by ids
// @Description Returns a list of loanReferralVisits that match the list of id.
// @Tags loanReferralVisits
// @Param data body types.ListLoanReferralVisitssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanReferralVisitssByIDsReply{}
// @Router /api/v1/loanReferralVisits/list/ids [post]
// @Security BearerAuth
func (h *loanReferralVisitsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanReferralVisitssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanReferralVisitsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanReferralVisitss := []*types.LoanReferralVisitsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanReferralVisitsMap[id]; ok {
			record, err := convertLoanReferralVisits(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanReferralVisits)
				return
			}
			loanReferralVisitss = append(loanReferralVisitss, record)
		}
	}

	response.Success(c, gin.H{
		"loanReferralVisitss": loanReferralVisitss,
	})
}

// ListByLastID get a paginated list of loanReferralVisitss by last id
// @Summary Get a paginated list of loanReferralVisitss by last id
// @Description Returns a paginated list of loanReferralVisitss starting after a given last id, useful for cursor-based pagination.
// @Tags loanReferralVisits
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanReferralVisitssReply{}
// @Router /api/v1/loanReferralVisits/list [get]
// @Security BearerAuth
func (h *loanReferralVisitsHandler) ListByLastID(c *gin.Context) {
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
	loanReferralVisitss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanReferralVisitss(loanReferralVisitss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanReferralVisits)
		return
	}

	response.Success(c, gin.H{
		"loanReferralVisitss": data,
	})
}

func getLoanReferralVisitsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanReferralVisits(loanReferralVisits *model.LoanReferralVisits) (*types.LoanReferralVisitsObjDetail, error) {
	data := &types.LoanReferralVisitsObjDetail{}
	err := copier.Copy(data, loanReferralVisits)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanReferralVisitss(fromValues []*model.LoanReferralVisits) ([]*types.LoanReferralVisitsObjDetail, error) {
	toValues := []*types.LoanReferralVisitsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanReferralVisits(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
