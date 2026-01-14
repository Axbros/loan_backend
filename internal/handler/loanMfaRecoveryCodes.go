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

var _ LoanMfaRecoveryCodesHandler = (*loanMfaRecoveryCodesHandler)(nil)

// LoanMfaRecoveryCodesHandler defining the handler interface
type LoanMfaRecoveryCodesHandler interface {
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

type loanMfaRecoveryCodesHandler struct {
	iDao dao.LoanMfaRecoveryCodesDao
}

// NewLoanMfaRecoveryCodesHandler creating the handler interface
func NewLoanMfaRecoveryCodesHandler() LoanMfaRecoveryCodesHandler {
	return &loanMfaRecoveryCodesHandler{
		iDao: dao.NewLoanMfaRecoveryCodesDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanMfaRecoveryCodesCache(database.GetCacheType()),
		),
	}
}

// Create a new loanMfaRecoveryCodes
// @Summary Create a new loanMfaRecoveryCodes
// @Description Creates a new loanMfaRecoveryCodes entity using the provided data in the request body.
// @Tags loanMfaRecoveryCodes
// @Accept json
// @Produce json
// @Param data body types.CreateLoanMfaRecoveryCodesRequest true "loanMfaRecoveryCodes information"
// @Success 200 {object} types.CreateLoanMfaRecoveryCodesReply{}
// @Router /api/v1/loanMfaRecoveryCodes [post]
// @Security BearerAuth
func (h *loanMfaRecoveryCodesHandler) Create(c *gin.Context) {
	form := &types.CreateLoanMfaRecoveryCodesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanMfaRecoveryCodes := &model.LoanMfaRecoveryCodes{}
	err = copier.Copy(loanMfaRecoveryCodes, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanMfaRecoveryCodes)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanMfaRecoveryCodes)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanMfaRecoveryCodes.ID})
}

// DeleteByID delete a loanMfaRecoveryCodes by id
// @Summary Delete a loanMfaRecoveryCodes by id
// @Description Deletes a existing loanMfaRecoveryCodes identified by the given id in the path.
// @Tags loanMfaRecoveryCodes
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanMfaRecoveryCodesByIDReply{}
// @Router /api/v1/loanMfaRecoveryCodes/{id} [delete]
// @Security BearerAuth
func (h *loanMfaRecoveryCodesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanMfaRecoveryCodesIDFromPath(c)
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

// UpdateByID update a loanMfaRecoveryCodes by id
// @Summary Update a loanMfaRecoveryCodes by id
// @Description Updates the specified loanMfaRecoveryCodes by given id in the path, support partial update.
// @Tags loanMfaRecoveryCodes
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanMfaRecoveryCodesByIDRequest true "loanMfaRecoveryCodes information"
// @Success 200 {object} types.UpdateLoanMfaRecoveryCodesByIDReply{}
// @Router /api/v1/loanMfaRecoveryCodes/{id} [put]
// @Security BearerAuth
func (h *loanMfaRecoveryCodesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanMfaRecoveryCodesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanMfaRecoveryCodesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanMfaRecoveryCodes := &model.LoanMfaRecoveryCodes{}
	err = copier.Copy(loanMfaRecoveryCodes, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanMfaRecoveryCodes)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanMfaRecoveryCodes)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanMfaRecoveryCodes by id
// @Summary Get a loanMfaRecoveryCodes by id
// @Description Gets detailed information of a loanMfaRecoveryCodes specified by the given id in the path.
// @Tags loanMfaRecoveryCodes
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanMfaRecoveryCodesByIDReply{}
// @Router /api/v1/loanMfaRecoveryCodes/{id} [get]
// @Security BearerAuth
func (h *loanMfaRecoveryCodesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanMfaRecoveryCodesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanMfaRecoveryCodes, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanMfaRecoveryCodesObjDetail{}
	err = copier.Copy(data, loanMfaRecoveryCodes)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanMfaRecoveryCodes)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanMfaRecoveryCodes": data})
}

// List get a paginated list of loanMfaRecoveryCodess by custom conditions
// @Summary Get a paginated list of loanMfaRecoveryCodess by custom conditions
// @Description Returns a paginated list of loanMfaRecoveryCodes based on query filters, including page number and size.
// @Tags loanMfaRecoveryCodes
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanMfaRecoveryCodessReply{}
// @Router /api/v1/loanMfaRecoveryCodes/list [post]
// @Security BearerAuth
func (h *loanMfaRecoveryCodesHandler) List(c *gin.Context) {
	form := &types.ListLoanMfaRecoveryCodessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanMfaRecoveryCodess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanMfaRecoveryCodess(loanMfaRecoveryCodess)
	if err != nil {
		response.Error(c, ecode.ErrListLoanMfaRecoveryCodes)
		return
	}

	response.Success(c, gin.H{
		"loanMfaRecoveryCodess": data,
		"total":                 total,
	})
}

// DeleteByIDs batch delete loanMfaRecoveryCodes by ids
// @Summary Batch delete loanMfaRecoveryCodes by ids
// @Description Deletes multiple loanMfaRecoveryCodes by a list of id
// @Tags loanMfaRecoveryCodes
// @Param data body types.DeleteLoanMfaRecoveryCodessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanMfaRecoveryCodessByIDsReply{}
// @Router /api/v1/loanMfaRecoveryCodes/delete/ids [post]
// @Security BearerAuth
func (h *loanMfaRecoveryCodesHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanMfaRecoveryCodessByIDsRequest{}
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

// GetByCondition get a loanMfaRecoveryCodes by custom condition
// @Summary Get a loanMfaRecoveryCodes by custom condition
// @Description Returns a single loanMfaRecoveryCodes that matches the specified filter conditions.
// @Tags loanMfaRecoveryCodes
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanMfaRecoveryCodesByConditionReply{}
// @Router /api/v1/loanMfaRecoveryCodes/condition [post]
// @Security BearerAuth
func (h *loanMfaRecoveryCodesHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanMfaRecoveryCodesByConditionRequest{}
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
	loanMfaRecoveryCodes, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanMfaRecoveryCodesObjDetail{}
	err = copier.Copy(data, loanMfaRecoveryCodes)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanMfaRecoveryCodes)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanMfaRecoveryCodes": data})
}

// ListByIDs batch get loanMfaRecoveryCodes by ids
// @Summary Batch get loanMfaRecoveryCodes by ids
// @Description Returns a list of loanMfaRecoveryCodes that match the list of id.
// @Tags loanMfaRecoveryCodes
// @Param data body types.ListLoanMfaRecoveryCodessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanMfaRecoveryCodessByIDsReply{}
// @Router /api/v1/loanMfaRecoveryCodes/list/ids [post]
// @Security BearerAuth
func (h *loanMfaRecoveryCodesHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanMfaRecoveryCodessByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanMfaRecoveryCodesMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanMfaRecoveryCodess := []*types.LoanMfaRecoveryCodesObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanMfaRecoveryCodesMap[id]; ok {
			record, err := convertLoanMfaRecoveryCodes(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanMfaRecoveryCodes)
				return
			}
			loanMfaRecoveryCodess = append(loanMfaRecoveryCodess, record)
		}
	}

	response.Success(c, gin.H{
		"loanMfaRecoveryCodess": loanMfaRecoveryCodess,
	})
}

// ListByLastID get a paginated list of loanMfaRecoveryCodess by last id
// @Summary Get a paginated list of loanMfaRecoveryCodess by last id
// @Description Returns a paginated list of loanMfaRecoveryCodess starting after a given last id, useful for cursor-based pagination.
// @Tags loanMfaRecoveryCodes
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanMfaRecoveryCodessReply{}
// @Router /api/v1/loanMfaRecoveryCodes/list [get]
// @Security BearerAuth
func (h *loanMfaRecoveryCodesHandler) ListByLastID(c *gin.Context) {
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
	loanMfaRecoveryCodess, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanMfaRecoveryCodess(loanMfaRecoveryCodess)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanMfaRecoveryCodes)
		return
	}

	response.Success(c, gin.H{
		"loanMfaRecoveryCodess": data,
	})
}

func getLoanMfaRecoveryCodesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanMfaRecoveryCodes(loanMfaRecoveryCodes *model.LoanMfaRecoveryCodes) (*types.LoanMfaRecoveryCodesObjDetail, error) {
	data := &types.LoanMfaRecoveryCodesObjDetail{}
	err := copier.Copy(data, loanMfaRecoveryCodes)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanMfaRecoveryCodess(fromValues []*model.LoanMfaRecoveryCodes) ([]*types.LoanMfaRecoveryCodesObjDetail, error) {
	toValues := []*types.LoanMfaRecoveryCodesObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanMfaRecoveryCodes(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
