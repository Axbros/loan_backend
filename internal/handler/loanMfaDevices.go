package handler

import (
	"errors"

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
