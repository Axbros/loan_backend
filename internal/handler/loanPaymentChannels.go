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

var _ LoanPaymentChannelsHandler = (*loanPaymentChannelsHandler)(nil)

// LoanPaymentChannelsHandler defining the handler interface
type LoanPaymentChannelsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type loanPaymentChannelsHandler struct {
	iDao dao.LoanPaymentChannelsDao
}

// NewLoanPaymentChannelsHandler creating the handler interface
func NewLoanPaymentChannelsHandler() LoanPaymentChannelsHandler {
	return &loanPaymentChannelsHandler{
		iDao: dao.NewLoanPaymentChannelsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanPaymentChannelsCache(database.GetCacheType()),
		),
	}
}

// Create a new loanPaymentChannels
// @Summary Create a new loanPaymentChannels
// @Description Creates a new loanPaymentChannels entity using the provided data in the request body.
// @Tags loanPaymentChannels
// @Accept json
// @Produce json
// @Param data body types.CreateLoanPaymentChannelsRequest true "loanPaymentChannels information"
// @Success 200 {object} types.CreateLoanPaymentChannelsReply{}
// @Router /api/v1/loanPaymentChannels [post]
// @Security BearerAuth
func (h *loanPaymentChannelsHandler) Create(c *gin.Context) {
	form := &types.CreateLoanPaymentChannelsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanPaymentChannels := &model.LoanPaymentChannels{}
	err = copier.Copy(loanPaymentChannels, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanPaymentChannels)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanPaymentChannels)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanPaymentChannels.ID})
}

// DeleteByID delete a loanPaymentChannels by id
// @Summary Delete a loanPaymentChannels by id
// @Description Deletes a existing loanPaymentChannels identified by the given id in the path.
// @Tags loanPaymentChannels
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanPaymentChannelsByIDReply{}
// @Router /api/v1/loanPaymentChannels/{id} [delete]
// @Security BearerAuth
func (h *loanPaymentChannelsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanPaymentChannelsIDFromPath(c)
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

// UpdateByID update a loanPaymentChannels by id
// @Summary Update a loanPaymentChannels by id
// @Description Updates the specified loanPaymentChannels by given id in the path, support partial update.
// @Tags loanPaymentChannels
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanPaymentChannelsByIDRequest true "loanPaymentChannels information"
// @Success 200 {object} types.UpdateLoanPaymentChannelsByIDReply{}
// @Router /api/v1/loanPaymentChannels/{id} [put]
// @Security BearerAuth
func (h *loanPaymentChannelsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanPaymentChannelsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanPaymentChannelsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanPaymentChannels := &model.LoanPaymentChannels{}
	err = copier.Copy(loanPaymentChannels, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanPaymentChannels)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanPaymentChannels)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanPaymentChannels by id
// @Summary Get a loanPaymentChannels by id
// @Description Gets detailed information of a loanPaymentChannels specified by the given id in the path.
// @Tags loanPaymentChannels
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanPaymentChannelsByIDReply{}
// @Router /api/v1/loanPaymentChannels/{id} [get]
// @Security BearerAuth
func (h *loanPaymentChannelsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanPaymentChannelsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanPaymentChannels, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanPaymentChannelsObjDetail{}
	err = copier.Copy(data, loanPaymentChannels)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanPaymentChannels)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanPaymentChannels": data})
}

// List get a paginated list of loanPaymentChannelss by custom conditions
// @Summary Get a paginated list of loanPaymentChannelss by custom conditions
// @Description Returns a paginated list of loanPaymentChannels based on query filters, including page number and size.
// @Tags loanPaymentChannels
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanPaymentChannelssReply{}
// @Router /api/v1/loanPaymentChannels/list [post]
// @Security BearerAuth
func (h *loanPaymentChannelsHandler) List(c *gin.Context) {
	form := &types.ListLoanPaymentChannelssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanPaymentChannelss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanPaymentChannelss(loanPaymentChannelss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanPaymentChannels)
		return
	}

	response.Success(c, gin.H{
		"loanPaymentChannelss": data,
		"total":                total,
	})
}

func getLoanPaymentChannelsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanPaymentChannels(loanPaymentChannels *model.LoanPaymentChannels) (*types.LoanPaymentChannelsObjDetail, error) {
	data := &types.LoanPaymentChannelsObjDetail{}
	err := copier.Copy(data, loanPaymentChannels)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanPaymentChannelss(fromValues []*model.LoanPaymentChannels) ([]*types.LoanPaymentChannelsObjDetail, error) {
	toValues := []*types.LoanPaymentChannelsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanPaymentChannels(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
