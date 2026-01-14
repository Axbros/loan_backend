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

var _ LoanPermissionsHandler = (*loanPermissionsHandler)(nil)

// LoanPermissionsHandler defining the handler interface
type LoanPermissionsHandler interface {
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

type loanPermissionsHandler struct {
	iDao dao.LoanPermissionsDao
}

// NewLoanPermissionsHandler creating the handler interface
func NewLoanPermissionsHandler() LoanPermissionsHandler {
	return &loanPermissionsHandler{
		iDao: dao.NewLoanPermissionsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanPermissionsCache(database.GetCacheType()),
		),
	}
}

// Create a new loanPermissions
// @Summary Create a new loanPermissions
// @Description Creates a new loanPermissions entity using the provided data in the request body.
// @Tags loanPermissions
// @Accept json
// @Produce json
// @Param data body types.CreateLoanPermissionsRequest true "loanPermissions information"
// @Success 200 {object} types.CreateLoanPermissionsReply{}
// @Router /api/v1/loanPermissions [post]
// @Security BearerAuth
func (h *loanPermissionsHandler) Create(c *gin.Context) {
	form := &types.CreateLoanPermissionsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanPermissions := &model.LoanPermissions{}
	err = copier.Copy(loanPermissions, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanPermissions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanPermissions)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanPermissions.ID})
}

// DeleteByID delete a loanPermissions by id
// @Summary Delete a loanPermissions by id
// @Description Deletes a existing loanPermissions identified by the given id in the path.
// @Tags loanPermissions
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanPermissionsByIDReply{}
// @Router /api/v1/loanPermissions/{id} [delete]
// @Security BearerAuth
func (h *loanPermissionsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanPermissionsIDFromPath(c)
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

// UpdateByID update a loanPermissions by id
// @Summary Update a loanPermissions by id
// @Description Updates the specified loanPermissions by given id in the path, support partial update.
// @Tags loanPermissions
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanPermissionsByIDRequest true "loanPermissions information"
// @Success 200 {object} types.UpdateLoanPermissionsByIDReply{}
// @Router /api/v1/loanPermissions/{id} [put]
// @Security BearerAuth
func (h *loanPermissionsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanPermissionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanPermissionsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanPermissions := &model.LoanPermissions{}
	err = copier.Copy(loanPermissions, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanPermissions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanPermissions)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanPermissions by id
// @Summary Get a loanPermissions by id
// @Description Gets detailed information of a loanPermissions specified by the given id in the path.
// @Tags loanPermissions
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanPermissionsByIDReply{}
// @Router /api/v1/loanPermissions/{id} [get]
// @Security BearerAuth
func (h *loanPermissionsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanPermissionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanPermissions, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanPermissionsObjDetail{}
	err = copier.Copy(data, loanPermissions)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanPermissions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanPermissions": data})
}

// List get a paginated list of loanPermissionss by custom conditions
// @Summary Get a paginated list of loanPermissionss by custom conditions
// @Description Returns a paginated list of loanPermissions based on query filters, including page number and size.
// @Tags loanPermissions
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanPermissionssReply{}
// @Router /api/v1/loanPermissions/list [post]
// @Security BearerAuth
func (h *loanPermissionsHandler) List(c *gin.Context) {
	form := &types.ListLoanPermissionssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanPermissionss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanPermissionss(loanPermissionss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanPermissions)
		return
	}

	response.Success(c, gin.H{
		"loanPermissionss": data,
		"total":            total,
	})
}

// DeleteByIDs batch delete loanPermissions by ids
// @Summary Batch delete loanPermissions by ids
// @Description Deletes multiple loanPermissions by a list of id
// @Tags loanPermissions
// @Param data body types.DeleteLoanPermissionssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanPermissionssByIDsReply{}
// @Router /api/v1/loanPermissions/delete/ids [post]
// @Security BearerAuth
func (h *loanPermissionsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanPermissionssByIDsRequest{}
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

// GetByCondition get a loanPermissions by custom condition
// @Summary Get a loanPermissions by custom condition
// @Description Returns a single loanPermissions that matches the specified filter conditions.
// @Tags loanPermissions
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanPermissionsByConditionReply{}
// @Router /api/v1/loanPermissions/condition [post]
// @Security BearerAuth
func (h *loanPermissionsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanPermissionsByConditionRequest{}
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
	loanPermissions, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanPermissionsObjDetail{}
	err = copier.Copy(data, loanPermissions)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanPermissions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanPermissions": data})
}

// ListByIDs batch get loanPermissions by ids
// @Summary Batch get loanPermissions by ids
// @Description Returns a list of loanPermissions that match the list of id.
// @Tags loanPermissions
// @Param data body types.ListLoanPermissionssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanPermissionssByIDsReply{}
// @Router /api/v1/loanPermissions/list/ids [post]
// @Security BearerAuth
func (h *loanPermissionsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanPermissionssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanPermissionsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanPermissionss := []*types.LoanPermissionsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanPermissionsMap[id]; ok {
			record, err := convertLoanPermissions(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanPermissions)
				return
			}
			loanPermissionss = append(loanPermissionss, record)
		}
	}

	response.Success(c, gin.H{
		"loanPermissionss": loanPermissionss,
	})
}

// ListByLastID get a paginated list of loanPermissionss by last id
// @Summary Get a paginated list of loanPermissionss by last id
// @Description Returns a paginated list of loanPermissionss starting after a given last id, useful for cursor-based pagination.
// @Tags loanPermissions
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanPermissionssReply{}
// @Router /api/v1/loanPermissions/list [get]
// @Security BearerAuth
func (h *loanPermissionsHandler) ListByLastID(c *gin.Context) {
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
	loanPermissionss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanPermissionss(loanPermissionss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanPermissions)
		return
	}

	response.Success(c, gin.H{
		"loanPermissionss": data,
	})
}

func getLoanPermissionsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanPermissions(loanPermissions *model.LoanPermissions) (*types.LoanPermissionsObjDetail, error) {
	data := &types.LoanPermissionsObjDetail{}
	err := copier.Copy(data, loanPermissions)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanPermissionss(fromValues []*model.LoanPermissions) ([]*types.LoanPermissionsObjDetail, error) {
	toValues := []*types.LoanPermissionsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanPermissions(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
