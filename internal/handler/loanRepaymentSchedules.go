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

var _ LoanRepaymentSchedulesHandler = (*loanRepaymentSchedulesHandler)(nil)

// LoanRepaymentSchedulesHandler defining the handler interface
type LoanRepaymentSchedulesHandler interface {
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

type loanRepaymentSchedulesHandler struct {
	iDao dao.LoanRepaymentSchedulesDao
}

// NewLoanRepaymentSchedulesHandler creating the handler interface
func NewLoanRepaymentSchedulesHandler() LoanRepaymentSchedulesHandler {
	return &loanRepaymentSchedulesHandler{
		iDao: dao.NewLoanRepaymentSchedulesDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanRepaymentSchedulesCache(database.GetCacheType()),
		),
	}
}

// Create a new loanRepaymentSchedules
// @Summary Create a new loanRepaymentSchedules
// @Description Creates a new loanRepaymentSchedules entity using the provided data in the request body.
// @Tags loanRepaymentSchedules
// @Accept json
// @Produce json
// @Param data body types.CreateLoanRepaymentSchedulesRequest true "loanRepaymentSchedules information"
// @Success 200 {object} types.CreateLoanRepaymentSchedulesReply{}
// @Router /api/v1/loanRepaymentSchedules [post]
// @Security BearerAuth
func (h *loanRepaymentSchedulesHandler) Create(c *gin.Context) {
	form := &types.CreateLoanRepaymentSchedulesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanRepaymentSchedules := &model.LoanRepaymentSchedules{}
	err = copier.Copy(loanRepaymentSchedules, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanRepaymentSchedules)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanRepaymentSchedules)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanRepaymentSchedules.ID})
}

// DeleteByID delete a loanRepaymentSchedules by id
// @Summary Delete a loanRepaymentSchedules by id
// @Description Deletes a existing loanRepaymentSchedules identified by the given id in the path.
// @Tags loanRepaymentSchedules
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanRepaymentSchedulesByIDReply{}
// @Router /api/v1/loanRepaymentSchedules/{id} [delete]
// @Security BearerAuth
func (h *loanRepaymentSchedulesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanRepaymentSchedulesIDFromPath(c)
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

// UpdateByID update a loanRepaymentSchedules by id
// @Summary Update a loanRepaymentSchedules by id
// @Description Updates the specified loanRepaymentSchedules by given id in the path, support partial update.
// @Tags loanRepaymentSchedules
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanRepaymentSchedulesByIDRequest true "loanRepaymentSchedules information"
// @Success 200 {object} types.UpdateLoanRepaymentSchedulesByIDReply{}
// @Router /api/v1/loanRepaymentSchedules/{id} [put]
// @Security BearerAuth
func (h *loanRepaymentSchedulesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanRepaymentSchedulesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanRepaymentSchedulesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanRepaymentSchedules := &model.LoanRepaymentSchedules{}
	err = copier.Copy(loanRepaymentSchedules, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanRepaymentSchedules)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanRepaymentSchedules)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanRepaymentSchedules by id
// @Summary Get a loanRepaymentSchedules by id
// @Description Gets detailed information of a loanRepaymentSchedules specified by the given id in the path.
// @Tags loanRepaymentSchedules
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanRepaymentSchedulesByIDReply{}
// @Router /api/v1/loanRepaymentSchedules/{id} [get]
// @Security BearerAuth
func (h *loanRepaymentSchedulesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanRepaymentSchedulesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRepaymentSchedules, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanRepaymentSchedulesObjDetail{}
	err = copier.Copy(data, loanRepaymentSchedules)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanRepaymentSchedules)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanRepaymentSchedules": data})
}

// List get a paginated list of loanRepaymentScheduless by custom conditions
// @Summary Get a paginated list of loanRepaymentScheduless by custom conditions
// @Description Returns a paginated list of loanRepaymentSchedules based on query filters, including page number and size.
// @Tags loanRepaymentSchedules
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanRepaymentSchedulessReply{}
// @Router /api/v1/loanRepaymentSchedules/list [post]
// @Security BearerAuth
func (h *loanRepaymentSchedulesHandler) List(c *gin.Context) {
	form := &types.ListLoanRepaymentSchedulessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRepaymentScheduless, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanRepaymentScheduless(loanRepaymentScheduless)
	if err != nil {
		response.Error(c, ecode.ErrListLoanRepaymentSchedules)
		return
	}

	response.Success(c, gin.H{
		"loanRepaymentScheduless": data,
		"total":                   total,
	})
}

// DeleteByIDs batch delete loanRepaymentSchedules by ids
// @Summary Batch delete loanRepaymentSchedules by ids
// @Description Deletes multiple loanRepaymentSchedules by a list of id
// @Tags loanRepaymentSchedules
// @Param data body types.DeleteLoanRepaymentSchedulessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanRepaymentSchedulessByIDsReply{}
// @Router /api/v1/loanRepaymentSchedules/delete/ids [post]
// @Security BearerAuth
func (h *loanRepaymentSchedulesHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanRepaymentSchedulessByIDsRequest{}
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

// GetByCondition get a loanRepaymentSchedules by custom condition
// @Summary Get a loanRepaymentSchedules by custom condition
// @Description Returns a single loanRepaymentSchedules that matches the specified filter conditions.
// @Tags loanRepaymentSchedules
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanRepaymentSchedulesByConditionReply{}
// @Router /api/v1/loanRepaymentSchedules/condition [post]
// @Security BearerAuth
func (h *loanRepaymentSchedulesHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanRepaymentSchedulesByConditionRequest{}
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
	loanRepaymentSchedules, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanRepaymentSchedulesObjDetail{}
	err = copier.Copy(data, loanRepaymentSchedules)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanRepaymentSchedules)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanRepaymentSchedules": data})
}

// ListByIDs batch get loanRepaymentSchedules by ids
// @Summary Batch get loanRepaymentSchedules by ids
// @Description Returns a list of loanRepaymentSchedules that match the list of id.
// @Tags loanRepaymentSchedules
// @Param data body types.ListLoanRepaymentSchedulessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanRepaymentSchedulessByIDsReply{}
// @Router /api/v1/loanRepaymentSchedules/list/ids [post]
// @Security BearerAuth
func (h *loanRepaymentSchedulesHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanRepaymentSchedulessByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRepaymentSchedulesMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanRepaymentScheduless := []*types.LoanRepaymentSchedulesObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanRepaymentSchedulesMap[id]; ok {
			record, err := convertLoanRepaymentSchedules(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanRepaymentSchedules)
				return
			}
			loanRepaymentScheduless = append(loanRepaymentScheduless, record)
		}
	}

	response.Success(c, gin.H{
		"loanRepaymentScheduless": loanRepaymentScheduless,
	})
}

// ListByLastID get a paginated list of loanRepaymentScheduless by last id
// @Summary Get a paginated list of loanRepaymentScheduless by last id
// @Description Returns a paginated list of loanRepaymentScheduless starting after a given last id, useful for cursor-based pagination.
// @Tags loanRepaymentSchedules
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanRepaymentSchedulessReply{}
// @Router /api/v1/loanRepaymentSchedules/list [get]
// @Security BearerAuth
func (h *loanRepaymentSchedulesHandler) ListByLastID(c *gin.Context) {
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
	loanRepaymentScheduless, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanRepaymentScheduless(loanRepaymentScheduless)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanRepaymentSchedules)
		return
	}

	response.Success(c, gin.H{
		"loanRepaymentScheduless": data,
	})
}

func getLoanRepaymentSchedulesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanRepaymentSchedules(loanRepaymentSchedules *model.LoanRepaymentSchedules) (*types.LoanRepaymentSchedulesObjDetail, error) {
	data := &types.LoanRepaymentSchedulesObjDetail{}
	err := copier.Copy(data, loanRepaymentSchedules)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanRepaymentScheduless(fromValues []*model.LoanRepaymentSchedules) ([]*types.LoanRepaymentSchedulesObjDetail, error) {
	toValues := []*types.LoanRepaymentSchedulesObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanRepaymentSchedules(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
