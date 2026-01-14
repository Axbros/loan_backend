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

var _ LoanPaymentChannelsHandler = (*loanPaymentChannelsHandler)(nil)

// LoanPaymentChannelsHandler defining the handler interface
type LoanPaymentChannelsHandler interface {
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

// DeleteByIDs batch delete loanPaymentChannels by ids
// @Summary Batch delete loanPaymentChannels by ids
// @Description Deletes multiple loanPaymentChannels by a list of id
// @Tags loanPaymentChannels
// @Param data body types.DeleteLoanPaymentChannelssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanPaymentChannelssByIDsReply{}
// @Router /api/v1/loanPaymentChannels/delete/ids [post]
// @Security BearerAuth
func (h *loanPaymentChannelsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanPaymentChannelssByIDsRequest{}
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

// GetByCondition get a loanPaymentChannels by custom condition
// @Summary Get a loanPaymentChannels by custom condition
// @Description Returns a single loanPaymentChannels that matches the specified filter conditions.
// @Tags loanPaymentChannels
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanPaymentChannelsByConditionReply{}
// @Router /api/v1/loanPaymentChannels/condition [post]
// @Security BearerAuth
func (h *loanPaymentChannelsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanPaymentChannelsByConditionRequest{}
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
	loanPaymentChannels, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanPaymentChannelsObjDetail{}
	err = copier.Copy(data, loanPaymentChannels)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanPaymentChannels)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanPaymentChannels": data})
}

// ListByIDs batch get loanPaymentChannels by ids
// @Summary Batch get loanPaymentChannels by ids
// @Description Returns a list of loanPaymentChannels that match the list of id.
// @Tags loanPaymentChannels
// @Param data body types.ListLoanPaymentChannelssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanPaymentChannelssByIDsReply{}
// @Router /api/v1/loanPaymentChannels/list/ids [post]
// @Security BearerAuth
func (h *loanPaymentChannelsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanPaymentChannelssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanPaymentChannelsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanPaymentChannelss := []*types.LoanPaymentChannelsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanPaymentChannelsMap[id]; ok {
			record, err := convertLoanPaymentChannels(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanPaymentChannels)
				return
			}
			loanPaymentChannelss = append(loanPaymentChannelss, record)
		}
	}

	response.Success(c, gin.H{
		"loanPaymentChannelss": loanPaymentChannelss,
	})
}

// ListByLastID get a paginated list of loanPaymentChannelss by last id
// @Summary Get a paginated list of loanPaymentChannelss by last id
// @Description Returns a paginated list of loanPaymentChannelss starting after a given last id, useful for cursor-based pagination.
// @Tags loanPaymentChannels
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanPaymentChannelssReply{}
// @Router /api/v1/loanPaymentChannels/list [get]
// @Security BearerAuth
func (h *loanPaymentChannelsHandler) ListByLastID(c *gin.Context) {
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
	loanPaymentChannelss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanPaymentChannelss(loanPaymentChannelss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanPaymentChannels)
		return
	}

	response.Success(c, gin.H{
		"loanPaymentChannelss": data,
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
