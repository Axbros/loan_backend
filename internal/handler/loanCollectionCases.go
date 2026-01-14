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

var _ LoanCollectionCasesHandler = (*loanCollectionCasesHandler)(nil)

// LoanCollectionCasesHandler defining the handler interface
type LoanCollectionCasesHandler interface {
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

type loanCollectionCasesHandler struct {
	iDao dao.LoanCollectionCasesDao
}

// NewLoanCollectionCasesHandler creating the handler interface
func NewLoanCollectionCasesHandler() LoanCollectionCasesHandler {
	return &loanCollectionCasesHandler{
		iDao: dao.NewLoanCollectionCasesDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanCollectionCasesCache(database.GetCacheType()),
		),
	}
}

// Create a new loanCollectionCases
// @Summary Create a new loanCollectionCases
// @Description Creates a new loanCollectionCases entity using the provided data in the request body.
// @Tags loanCollectionCases
// @Accept json
// @Produce json
// @Param data body types.CreateLoanCollectionCasesRequest true "loanCollectionCases information"
// @Success 200 {object} types.CreateLoanCollectionCasesReply{}
// @Router /api/v1/loanCollectionCases [post]
// @Security BearerAuth
func (h *loanCollectionCasesHandler) Create(c *gin.Context) {
	form := &types.CreateLoanCollectionCasesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanCollectionCases := &model.LoanCollectionCases{}
	err = copier.Copy(loanCollectionCases, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanCollectionCases)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanCollectionCases)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanCollectionCases.ID})
}

// DeleteByID delete a loanCollectionCases by id
// @Summary Delete a loanCollectionCases by id
// @Description Deletes a existing loanCollectionCases identified by the given id in the path.
// @Tags loanCollectionCases
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanCollectionCasesByIDReply{}
// @Router /api/v1/loanCollectionCases/{id} [delete]
// @Security BearerAuth
func (h *loanCollectionCasesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanCollectionCasesIDFromPath(c)
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

// UpdateByID update a loanCollectionCases by id
// @Summary Update a loanCollectionCases by id
// @Description Updates the specified loanCollectionCases by given id in the path, support partial update.
// @Tags loanCollectionCases
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanCollectionCasesByIDRequest true "loanCollectionCases information"
// @Success 200 {object} types.UpdateLoanCollectionCasesByIDReply{}
// @Router /api/v1/loanCollectionCases/{id} [put]
// @Security BearerAuth
func (h *loanCollectionCasesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanCollectionCasesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanCollectionCasesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanCollectionCases := &model.LoanCollectionCases{}
	err = copier.Copy(loanCollectionCases, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanCollectionCases)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanCollectionCases)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanCollectionCases by id
// @Summary Get a loanCollectionCases by id
// @Description Gets detailed information of a loanCollectionCases specified by the given id in the path.
// @Tags loanCollectionCases
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanCollectionCasesByIDReply{}
// @Router /api/v1/loanCollectionCases/{id} [get]
// @Security BearerAuth
func (h *loanCollectionCasesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanCollectionCasesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanCollectionCases, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanCollectionCasesObjDetail{}
	err = copier.Copy(data, loanCollectionCases)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanCollectionCases)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanCollectionCases": data})
}

// List get a paginated list of loanCollectionCasess by custom conditions
// @Summary Get a paginated list of loanCollectionCasess by custom conditions
// @Description Returns a paginated list of loanCollectionCases based on query filters, including page number and size.
// @Tags loanCollectionCases
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanCollectionCasessReply{}
// @Router /api/v1/loanCollectionCases/list [post]
// @Security BearerAuth
func (h *loanCollectionCasesHandler) List(c *gin.Context) {
	form := &types.ListLoanCollectionCasessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanCollectionCasess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanCollectionCasess(loanCollectionCasess)
	if err != nil {
		response.Error(c, ecode.ErrListLoanCollectionCases)
		return
	}

	response.Success(c, gin.H{
		"loanCollectionCasess": data,
		"total":                total,
	})
}

// DeleteByIDs batch delete loanCollectionCases by ids
// @Summary Batch delete loanCollectionCases by ids
// @Description Deletes multiple loanCollectionCases by a list of id
// @Tags loanCollectionCases
// @Param data body types.DeleteLoanCollectionCasessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanCollectionCasessByIDsReply{}
// @Router /api/v1/loanCollectionCases/delete/ids [post]
// @Security BearerAuth
func (h *loanCollectionCasesHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanCollectionCasessByIDsRequest{}
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

// GetByCondition get a loanCollectionCases by custom condition
// @Summary Get a loanCollectionCases by custom condition
// @Description Returns a single loanCollectionCases that matches the specified filter conditions.
// @Tags loanCollectionCases
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanCollectionCasesByConditionReply{}
// @Router /api/v1/loanCollectionCases/condition [post]
// @Security BearerAuth
func (h *loanCollectionCasesHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanCollectionCasesByConditionRequest{}
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
	loanCollectionCases, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanCollectionCasesObjDetail{}
	err = copier.Copy(data, loanCollectionCases)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanCollectionCases)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanCollectionCases": data})
}

// ListByIDs batch get loanCollectionCases by ids
// @Summary Batch get loanCollectionCases by ids
// @Description Returns a list of loanCollectionCases that match the list of id.
// @Tags loanCollectionCases
// @Param data body types.ListLoanCollectionCasessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanCollectionCasessByIDsReply{}
// @Router /api/v1/loanCollectionCases/list/ids [post]
// @Security BearerAuth
func (h *loanCollectionCasesHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanCollectionCasessByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanCollectionCasesMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanCollectionCasess := []*types.LoanCollectionCasesObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanCollectionCasesMap[id]; ok {
			record, err := convertLoanCollectionCases(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanCollectionCases)
				return
			}
			loanCollectionCasess = append(loanCollectionCasess, record)
		}
	}

	response.Success(c, gin.H{
		"loanCollectionCasess": loanCollectionCasess,
	})
}

// ListByLastID get a paginated list of loanCollectionCasess by last id
// @Summary Get a paginated list of loanCollectionCasess by last id
// @Description Returns a paginated list of loanCollectionCasess starting after a given last id, useful for cursor-based pagination.
// @Tags loanCollectionCases
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanCollectionCasessReply{}
// @Router /api/v1/loanCollectionCases/list [get]
// @Security BearerAuth
func (h *loanCollectionCasesHandler) ListByLastID(c *gin.Context) {
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
	loanCollectionCasess, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanCollectionCasess(loanCollectionCasess)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanCollectionCases)
		return
	}

	response.Success(c, gin.H{
		"loanCollectionCasess": data,
	})
}

func getLoanCollectionCasesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanCollectionCases(loanCollectionCases *model.LoanCollectionCases) (*types.LoanCollectionCasesObjDetail, error) {
	data := &types.LoanCollectionCasesObjDetail{}
	err := copier.Copy(data, loanCollectionCases)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanCollectionCasess(fromValues []*model.LoanCollectionCases) ([]*types.LoanCollectionCasesObjDetail, error) {
	toValues := []*types.LoanCollectionCasesObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanCollectionCases(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
