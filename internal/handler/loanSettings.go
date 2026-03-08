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

var _ LoanSettingsHandler = (*loanSettingsHandler)(nil)

// LoanSettingsHandler defining the handler interface
type LoanSettingsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type loanSettingsHandler struct {
	iDao dao.LoanSettingsDao
}

// NewLoanSettingsHandler creating the handler interface
func NewLoanSettingsHandler() LoanSettingsHandler {
	return &loanSettingsHandler{
		iDao: dao.NewLoanSettingsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanSettingsCache(database.GetCacheType()),
		),
	}
}

// Create a new loanSettings
// @Summary Create a new loanSettings
// @Description Creates a new loanSettings entity using the provided data in the request body.
// @Tags loanSettings
// @Accept json
// @Produce json
// @Param data body types.CreateLoanSettingsRequest true "loanSettings information"
// @Success 200 {object} types.CreateLoanSettingsReply{}
// @Router /api/v1/loanSettings [post]
// @Security BearerAuth
func (h *loanSettingsHandler) Create(c *gin.Context) {
	form := &types.CreateLoanSettingsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanSettings := &model.LoanSettings{}
	err = copier.Copy(loanSettings, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanSettings)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanSettings)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanSettings.ID})
}

// DeleteByID delete a loanSettings by id
// @Summary Delete a loanSettings by id
// @Description Deletes a existing loanSettings identified by the given id in the path.
// @Tags loanSettings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanSettingsByIDReply{}
// @Router /api/v1/loanSettings/{id} [delete]
// @Security BearerAuth
func (h *loanSettingsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanSettingsIDFromPath(c)
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

// UpdateByID update a loanSettings by id
// @Summary Update a loanSettings by id
// @Description Updates the specified loanSettings by given id in the path, support partial update.
// @Tags loanSettings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanSettingsByIDRequest true "loanSettings information"
// @Success 200 {object} types.UpdateLoanSettingsByIDReply{}
// @Router /api/v1/loanSettings/{id} [put]
// @Security BearerAuth
func (h *loanSettingsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanSettingsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanSettingsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanSettings := &model.LoanSettings{}
	err = copier.Copy(loanSettings, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanSettings)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanSettings)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanSettings by id
// @Summary Get a loanSettings by id
// @Description Gets detailed information of a loanSettings specified by the given id in the path.
// @Tags loanSettings
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanSettingsByIDReply{}
// @Router /api/v1/loanSettings/{id} [get]
// @Security BearerAuth
func (h *loanSettingsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanSettingsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanSettings, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanSettingsObjDetail{}
	err = copier.Copy(data, loanSettings)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanSettings)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanSettings": data})
}

// List get a paginated list of loanSettingss by custom conditions
// @Summary Get a paginated list of loanSettingss by custom conditions
// @Description Returns a paginated list of loanSettings based on query filters, including page number and size.
// @Tags loanSettings
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanSettingssReply{}
// @Router /api/v1/loanSettings/list [post]
// @Security BearerAuth
func (h *loanSettingsHandler) List(c *gin.Context) {
	form := &types.ListLoanSettingssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanSettingss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanSettingss(loanSettingss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanSettings)
		return
	}

	response.Success(c, gin.H{
		"records": data,
		"total":   total,
	})
}

func getLoanSettingsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanSettings(loanSettings *model.LoanSettings) (*types.LoanSettingsObjDetail, error) {
	data := &types.LoanSettingsObjDetail{}
	err := copier.Copy(data, loanSettings)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanSettingss(fromValues []*model.LoanSettings) ([]*types.LoanSettingsObjDetail, error) {
	toValues := []*types.LoanSettingsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanSettings(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
