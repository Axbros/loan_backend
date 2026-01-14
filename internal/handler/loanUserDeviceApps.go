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

var _ LoanUserDeviceAppsHandler = (*loanUserDeviceAppsHandler)(nil)

// LoanUserDeviceAppsHandler defining the handler interface
type LoanUserDeviceAppsHandler interface {
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

// DeleteByIDs batch delete loanUserDeviceApps by ids
// @Summary Batch delete loanUserDeviceApps by ids
// @Description Deletes multiple loanUserDeviceApps by a list of id
// @Tags loanUserDeviceApps
// @Param data body types.DeleteLoanUserDeviceAppssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanUserDeviceAppssByIDsReply{}
// @Router /api/v1/loanUserDeviceApps/delete/ids [post]
// @Security BearerAuth
func (h *loanUserDeviceAppsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanUserDeviceAppssByIDsRequest{}
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

// GetByCondition get a loanUserDeviceApps by custom condition
// @Summary Get a loanUserDeviceApps by custom condition
// @Description Returns a single loanUserDeviceApps that matches the specified filter conditions.
// @Tags loanUserDeviceApps
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanUserDeviceAppsByConditionReply{}
// @Router /api/v1/loanUserDeviceApps/condition [post]
// @Security BearerAuth
func (h *loanUserDeviceAppsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanUserDeviceAppsByConditionRequest{}
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
	loanUserDeviceApps, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanUserDeviceAppsObjDetail{}
	err = copier.Copy(data, loanUserDeviceApps)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanUserDeviceApps)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanUserDeviceApps": data})
}

// ListByIDs batch get loanUserDeviceApps by ids
// @Summary Batch get loanUserDeviceApps by ids
// @Description Returns a list of loanUserDeviceApps that match the list of id.
// @Tags loanUserDeviceApps
// @Param data body types.ListLoanUserDeviceAppssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanUserDeviceAppssByIDsReply{}
// @Router /api/v1/loanUserDeviceApps/list/ids [post]
// @Security BearerAuth
func (h *loanUserDeviceAppsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanUserDeviceAppssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUserDeviceAppsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanUserDeviceAppss := []*types.LoanUserDeviceAppsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanUserDeviceAppsMap[id]; ok {
			record, err := convertLoanUserDeviceApps(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanUserDeviceApps)
				return
			}
			loanUserDeviceAppss = append(loanUserDeviceAppss, record)
		}
	}

	response.Success(c, gin.H{
		"loanUserDeviceAppss": loanUserDeviceAppss,
	})
}

// ListByLastID get a paginated list of loanUserDeviceAppss by last id
// @Summary Get a paginated list of loanUserDeviceAppss by last id
// @Description Returns a paginated list of loanUserDeviceAppss starting after a given last id, useful for cursor-based pagination.
// @Tags loanUserDeviceApps
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanUserDeviceAppssReply{}
// @Router /api/v1/loanUserDeviceApps/list [get]
// @Security BearerAuth
func (h *loanUserDeviceAppsHandler) ListByLastID(c *gin.Context) {
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
	loanUserDeviceAppss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanUserDeviceAppss(loanUserDeviceAppss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanUserDeviceApps)
		return
	}

	response.Success(c, gin.H{
		"loanUserDeviceAppss": data,
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
