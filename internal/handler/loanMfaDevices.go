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

var _ LoanMfaDevicesHandler = (*loanMfaDevicesHandler)(nil)

// LoanMfaDevicesHandler defining the handler interface
type LoanMfaDevicesHandler interface {
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

type loanMfaDevicesHandler struct {
	iDao dao.LoanMfaDevicesDao
}

// NewLoanMfaDevicesHandler creating the handler interface
func NewLoanMfaDevicesHandler() LoanMfaDevicesHandler {
	return &loanMfaDevicesHandler{
		iDao: dao.NewLoanMfaDevicesDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanMfaDevicesCache(database.GetCacheType()),
		),
	}
}

// Create a new loanMfaDevices
// @Summary Create a new loanMfaDevices
// @Description Creates a new loanMfaDevices entity using the provided data in the request body.
// @Tags loanMfaDevices
// @Accept json
// @Produce json
// @Param data body types.CreateLoanMfaDevicesRequest true "loanMfaDevices information"
// @Success 200 {object} types.CreateLoanMfaDevicesReply{}
// @Router /api/v1/loanMfaDevices [post]
// @Security BearerAuth
func (h *loanMfaDevicesHandler) Create(c *gin.Context) {
	form := &types.CreateLoanMfaDevicesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanMfaDevices := &model.LoanMfaDevices{}
	err = copier.Copy(loanMfaDevices, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanMfaDevices)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanMfaDevices)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanMfaDevices.ID})
}

// DeleteByID delete a loanMfaDevices by id
// @Summary Delete a loanMfaDevices by id
// @Description Deletes a existing loanMfaDevices identified by the given id in the path.
// @Tags loanMfaDevices
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanMfaDevicesByIDReply{}
// @Router /api/v1/loanMfaDevices/{id} [delete]
// @Security BearerAuth
func (h *loanMfaDevicesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanMfaDevicesIDFromPath(c)
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

// UpdateByID update a loanMfaDevices by id
// @Summary Update a loanMfaDevices by id
// @Description Updates the specified loanMfaDevices by given id in the path, support partial update.
// @Tags loanMfaDevices
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanMfaDevicesByIDRequest true "loanMfaDevices information"
// @Success 200 {object} types.UpdateLoanMfaDevicesByIDReply{}
// @Router /api/v1/loanMfaDevices/{id} [put]
// @Security BearerAuth
func (h *loanMfaDevicesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanMfaDevicesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanMfaDevicesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanMfaDevices := &model.LoanMfaDevices{}
	err = copier.Copy(loanMfaDevices, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanMfaDevices)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanMfaDevices)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanMfaDevices by id
// @Summary Get a loanMfaDevices by id
// @Description Gets detailed information of a loanMfaDevices specified by the given id in the path.
// @Tags loanMfaDevices
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanMfaDevicesByIDReply{}
// @Router /api/v1/loanMfaDevices/{id} [get]
// @Security BearerAuth
func (h *loanMfaDevicesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanMfaDevicesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanMfaDevices, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanMfaDevicesObjDetail{}
	err = copier.Copy(data, loanMfaDevices)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanMfaDevices)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanMfaDevices": data})
}

// List get a paginated list of loanMfaDevicess by custom conditions
// @Summary Get a paginated list of loanMfaDevicess by custom conditions
// @Description Returns a paginated list of loanMfaDevices based on query filters, including page number and size.
// @Tags loanMfaDevices
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanMfaDevicessReply{}
// @Router /api/v1/loanMfaDevices/list [post]
// @Security BearerAuth
func (h *loanMfaDevicesHandler) List(c *gin.Context) {
	form := &types.ListLoanMfaDevicessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanMfaDevicess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanMfaDevicess(loanMfaDevicess)
	if err != nil {
		response.Error(c, ecode.ErrListLoanMfaDevices)
		return
	}

	response.Success(c, gin.H{
		"loanMfaDevicess": data,
		"total":           total,
	})
}

// DeleteByIDs batch delete loanMfaDevices by ids
// @Summary Batch delete loanMfaDevices by ids
// @Description Deletes multiple loanMfaDevices by a list of id
// @Tags loanMfaDevices
// @Param data body types.DeleteLoanMfaDevicessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanMfaDevicessByIDsReply{}
// @Router /api/v1/loanMfaDevices/delete/ids [post]
// @Security BearerAuth
func (h *loanMfaDevicesHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanMfaDevicessByIDsRequest{}
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

// GetByCondition get a loanMfaDevices by custom condition
// @Summary Get a loanMfaDevices by custom condition
// @Description Returns a single loanMfaDevices that matches the specified filter conditions.
// @Tags loanMfaDevices
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanMfaDevicesByConditionReply{}
// @Router /api/v1/loanMfaDevices/condition [post]
// @Security BearerAuth
func (h *loanMfaDevicesHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanMfaDevicesByConditionRequest{}
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
	loanMfaDevices, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanMfaDevicesObjDetail{}
	err = copier.Copy(data, loanMfaDevices)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanMfaDevices)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanMfaDevices": data})
}

// ListByIDs batch get loanMfaDevices by ids
// @Summary Batch get loanMfaDevices by ids
// @Description Returns a list of loanMfaDevices that match the list of id.
// @Tags loanMfaDevices
// @Param data body types.ListLoanMfaDevicessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanMfaDevicessByIDsReply{}
// @Router /api/v1/loanMfaDevices/list/ids [post]
// @Security BearerAuth
func (h *loanMfaDevicesHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanMfaDevicessByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanMfaDevicesMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanMfaDevicess := []*types.LoanMfaDevicesObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanMfaDevicesMap[id]; ok {
			record, err := convertLoanMfaDevices(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanMfaDevices)
				return
			}
			loanMfaDevicess = append(loanMfaDevicess, record)
		}
	}

	response.Success(c, gin.H{
		"loanMfaDevicess": loanMfaDevicess,
	})
}

// ListByLastID get a paginated list of loanMfaDevicess by last id
// @Summary Get a paginated list of loanMfaDevicess by last id
// @Description Returns a paginated list of loanMfaDevicess starting after a given last id, useful for cursor-based pagination.
// @Tags loanMfaDevices
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanMfaDevicessReply{}
// @Router /api/v1/loanMfaDevices/list [get]
// @Security BearerAuth
func (h *loanMfaDevicesHandler) ListByLastID(c *gin.Context) {
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
	loanMfaDevicess, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanMfaDevicess(loanMfaDevicess)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanMfaDevices)
		return
	}

	response.Success(c, gin.H{
		"loanMfaDevicess": data,
	})
}

func getLoanMfaDevicesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanMfaDevices(loanMfaDevices *model.LoanMfaDevices) (*types.LoanMfaDevicesObjDetail, error) {
	data := &types.LoanMfaDevicesObjDetail{}
	err := copier.Copy(data, loanMfaDevices)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanMfaDevicess(fromValues []*model.LoanMfaDevices) ([]*types.LoanMfaDevicesObjDetail, error) {
	toValues := []*types.LoanMfaDevicesObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanMfaDevices(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
