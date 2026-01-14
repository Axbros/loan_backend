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

var _ LoanRepaymentTransactionsHandler = (*loanRepaymentTransactionsHandler)(nil)

// LoanRepaymentTransactionsHandler defining the handler interface
type LoanRepaymentTransactionsHandler interface {
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

type loanRepaymentTransactionsHandler struct {
	iDao dao.LoanRepaymentTransactionsDao
}

// NewLoanRepaymentTransactionsHandler creating the handler interface
func NewLoanRepaymentTransactionsHandler() LoanRepaymentTransactionsHandler {
	return &loanRepaymentTransactionsHandler{
		iDao: dao.NewLoanRepaymentTransactionsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanRepaymentTransactionsCache(database.GetCacheType()),
		),
	}
}

// Create a new loanRepaymentTransactions
// @Summary Create a new loanRepaymentTransactions
// @Description Creates a new loanRepaymentTransactions entity using the provided data in the request body.
// @Tags loanRepaymentTransactions
// @Accept json
// @Produce json
// @Param data body types.CreateLoanRepaymentTransactionsRequest true "loanRepaymentTransactions information"
// @Success 200 {object} types.CreateLoanRepaymentTransactionsReply{}
// @Router /api/v1/loanRepaymentTransactions [post]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) Create(c *gin.Context) {
	form := &types.CreateLoanRepaymentTransactionsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanRepaymentTransactions := &model.LoanRepaymentTransactions{}
	err = copier.Copy(loanRepaymentTransactions, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanRepaymentTransactions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanRepaymentTransactions)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanRepaymentTransactions.ID})
}

// DeleteByID delete a loanRepaymentTransactions by id
// @Summary Delete a loanRepaymentTransactions by id
// @Description Deletes a existing loanRepaymentTransactions identified by the given id in the path.
// @Tags loanRepaymentTransactions
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanRepaymentTransactionsByIDReply{}
// @Router /api/v1/loanRepaymentTransactions/{id} [delete]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanRepaymentTransactionsIDFromPath(c)
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

// UpdateByID update a loanRepaymentTransactions by id
// @Summary Update a loanRepaymentTransactions by id
// @Description Updates the specified loanRepaymentTransactions by given id in the path, support partial update.
// @Tags loanRepaymentTransactions
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanRepaymentTransactionsByIDRequest true "loanRepaymentTransactions information"
// @Success 200 {object} types.UpdateLoanRepaymentTransactionsByIDReply{}
// @Router /api/v1/loanRepaymentTransactions/{id} [put]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanRepaymentTransactionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanRepaymentTransactionsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanRepaymentTransactions := &model.LoanRepaymentTransactions{}
	err = copier.Copy(loanRepaymentTransactions, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanRepaymentTransactions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanRepaymentTransactions)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanRepaymentTransactions by id
// @Summary Get a loanRepaymentTransactions by id
// @Description Gets detailed information of a loanRepaymentTransactions specified by the given id in the path.
// @Tags loanRepaymentTransactions
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanRepaymentTransactionsByIDReply{}
// @Router /api/v1/loanRepaymentTransactions/{id} [get]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanRepaymentTransactionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRepaymentTransactions, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanRepaymentTransactionsObjDetail{}
	err = copier.Copy(data, loanRepaymentTransactions)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanRepaymentTransactions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanRepaymentTransactions": data})
}

// List get a paginated list of loanRepaymentTransactionss by custom conditions
// @Summary Get a paginated list of loanRepaymentTransactionss by custom conditions
// @Description Returns a paginated list of loanRepaymentTransactions based on query filters, including page number and size.
// @Tags loanRepaymentTransactions
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanRepaymentTransactionssReply{}
// @Router /api/v1/loanRepaymentTransactions/list [post]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) List(c *gin.Context) {
	form := &types.ListLoanRepaymentTransactionssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRepaymentTransactionss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanRepaymentTransactionss(loanRepaymentTransactionss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanRepaymentTransactions)
		return
	}

	response.Success(c, gin.H{
		"loanRepaymentTransactionss": data,
		"total":                      total,
	})
}

// DeleteByIDs batch delete loanRepaymentTransactions by ids
// @Summary Batch delete loanRepaymentTransactions by ids
// @Description Deletes multiple loanRepaymentTransactions by a list of id
// @Tags loanRepaymentTransactions
// @Param data body types.DeleteLoanRepaymentTransactionssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanRepaymentTransactionssByIDsReply{}
// @Router /api/v1/loanRepaymentTransactions/delete/ids [post]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanRepaymentTransactionssByIDsRequest{}
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

// GetByCondition get a loanRepaymentTransactions by custom condition
// @Summary Get a loanRepaymentTransactions by custom condition
// @Description Returns a single loanRepaymentTransactions that matches the specified filter conditions.
// @Tags loanRepaymentTransactions
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanRepaymentTransactionsByConditionReply{}
// @Router /api/v1/loanRepaymentTransactions/condition [post]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanRepaymentTransactionsByConditionRequest{}
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
	loanRepaymentTransactions, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanRepaymentTransactionsObjDetail{}
	err = copier.Copy(data, loanRepaymentTransactions)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanRepaymentTransactions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanRepaymentTransactions": data})
}

// ListByIDs batch get loanRepaymentTransactions by ids
// @Summary Batch get loanRepaymentTransactions by ids
// @Description Returns a list of loanRepaymentTransactions that match the list of id.
// @Tags loanRepaymentTransactions
// @Param data body types.ListLoanRepaymentTransactionssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanRepaymentTransactionssByIDsReply{}
// @Router /api/v1/loanRepaymentTransactions/list/ids [post]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanRepaymentTransactionssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRepaymentTransactionsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanRepaymentTransactionss := []*types.LoanRepaymentTransactionsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanRepaymentTransactionsMap[id]; ok {
			record, err := convertLoanRepaymentTransactions(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanRepaymentTransactions)
				return
			}
			loanRepaymentTransactionss = append(loanRepaymentTransactionss, record)
		}
	}

	response.Success(c, gin.H{
		"loanRepaymentTransactionss": loanRepaymentTransactionss,
	})
}

// ListByLastID get a paginated list of loanRepaymentTransactionss by last id
// @Summary Get a paginated list of loanRepaymentTransactionss by last id
// @Description Returns a paginated list of loanRepaymentTransactionss starting after a given last id, useful for cursor-based pagination.
// @Tags loanRepaymentTransactions
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanRepaymentTransactionssReply{}
// @Router /api/v1/loanRepaymentTransactions/list [get]
// @Security BearerAuth
func (h *loanRepaymentTransactionsHandler) ListByLastID(c *gin.Context) {
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
	loanRepaymentTransactionss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanRepaymentTransactionss(loanRepaymentTransactionss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanRepaymentTransactions)
		return
	}

	response.Success(c, gin.H{
		"loanRepaymentTransactionss": data,
	})
}

func getLoanRepaymentTransactionsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanRepaymentTransactions(loanRepaymentTransactions *model.LoanRepaymentTransactions) (*types.LoanRepaymentTransactionsObjDetail, error) {
	data := &types.LoanRepaymentTransactionsObjDetail{}
	err := copier.Copy(data, loanRepaymentTransactions)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanRepaymentTransactionss(fromValues []*model.LoanRepaymentTransactions) ([]*types.LoanRepaymentTransactionsObjDetail, error) {
	toValues := []*types.LoanRepaymentTransactionsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanRepaymentTransactions(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
