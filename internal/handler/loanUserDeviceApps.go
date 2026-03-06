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

var _ LoanUserDeviceAppsHandler = (*loanUserDeviceAppsHandler)(nil)

// LoanUserDeviceAppsHandler defining the handler interface
type LoanUserDeviceAppsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type loanUserDeviceAppsHandler struct {
	iDao dao.LoanUserDeviceAppsDao
}

// NewLoanUserDeviceAppsHandler creating the handler interface
func NewLoanUserDeviceAppsHandler() LoanUserDeviceAppsHandler {
	return &loanUserDeviceAppsHandler{
		iDao: dao.NewLoanUserDeviceAppsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanUserDeviceAppsCache(database.GetCacheType()),
		),
	}
}

// Create a new loanUserDeviceApps
// @Summary Create a new loanUserDeviceApps
// @Description Creates a new loanUserDeviceApps entity using the provided data in the request body.
// @Tags loanUserDeviceApps
// @Accept json
// @Produce json
// @Param data body types.CreateLoanUserDeviceAppsRequest true "loanUserDeviceApps information"
// @Success 200 {object} types.CreateLoanUserDeviceAppsReply{}
// @Router /api/v1/loanUserDeviceApps [post]
// @Security BearerAuth
func (h *loanUserDeviceAppsHandler) Create(c *gin.Context) {
	form := &types.CreateLoanUserDeviceAppsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanUserDeviceApps := &model.LoanUserDeviceApps{}
	err = copier.Copy(loanUserDeviceApps, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanUserDeviceApps)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanUserDeviceApps)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanUserDeviceApps.ID})
}

// DeleteByID delete a loanUserDeviceApps by id
// @Summary Delete a loanUserDeviceApps by id
// @Description Deletes a existing loanUserDeviceApps identified by the given id in the path.
// @Tags loanUserDeviceApps
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanUserDeviceAppsByIDReply{}
// @Router /api/v1/loanUserDeviceApps/{id} [delete]
// @Security BearerAuth
func (h *loanUserDeviceAppsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanUserDeviceAppsIDFromPath(c)
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

// UpdateByID update a loanUserDeviceApps by id
// @Summary Update a loanUserDeviceApps by id
// @Description Updates the specified loanUserDeviceApps by given id in the path, support partial update.
// @Tags loanUserDeviceApps
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanUserDeviceAppsByIDRequest true "loanUserDeviceApps information"
// @Success 200 {object} types.UpdateLoanUserDeviceAppsByIDReply{}
// @Router /api/v1/loanUserDeviceApps/{id} [put]
// @Security BearerAuth
func (h *loanUserDeviceAppsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanUserDeviceAppsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanUserDeviceAppsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanUserDeviceApps := &model.LoanUserDeviceApps{}
	err = copier.Copy(loanUserDeviceApps, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanUserDeviceApps)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanUserDeviceApps)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanUserDeviceApps by id
// @Summary Get a loanUserDeviceApps by id
// @Description Gets detailed information of a loanUserDeviceApps specified by the given id in the path.
// @Tags loanUserDeviceApps
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanUserDeviceAppsByIDReply{}
// @Router /api/v1/loanUserDeviceApps/{id} [get]
// @Security BearerAuth
func (h *loanUserDeviceAppsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanUserDeviceAppsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUserDeviceApps, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanUserDeviceAppsObjDetail{}
	err = copier.Copy(data, loanUserDeviceApps)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanUserDeviceApps)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanUserDeviceApps": data})
}

// List get a paginated list of loanUserDeviceAppss by custom conditions
// @Summary Get a paginated list of loanUserDeviceAppss by custom conditions
// @Description Returns a paginated list of loanUserDeviceApps based on query filters, including page number and size.
// @Tags loanUserDeviceApps
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanUserDeviceAppssReply{}
// @Router /api/v1/loanUserDeviceApps/list [post]
// @Security BearerAuth
func (h *loanUserDeviceAppsHandler) List(c *gin.Context) {
	form := &types.ListLoanUserDeviceAppssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUserDeviceAppss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanUserDeviceAppss(loanUserDeviceAppss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanUserDeviceApps)
		return
	}

	response.Success(c, gin.H{
		"loanUserDeviceAppss": data,
		"total":               total,
	})
}

func getLoanUserDeviceAppsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanUserDeviceApps(loanUserDeviceApps *model.LoanUserDeviceApps) (*types.LoanUserDeviceAppsObjDetail, error) {
	data := &types.LoanUserDeviceAppsObjDetail{}
	err := copier.Copy(data, loanUserDeviceApps)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanUserDeviceAppss(fromValues []*model.LoanUserDeviceApps) ([]*types.LoanUserDeviceAppsObjDetail, error) {
	toValues := []*types.LoanUserDeviceAppsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanUserDeviceApps(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
