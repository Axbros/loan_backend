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

var _ LoanRolePermissionsHandler = (*loanRolePermissionsHandler)(nil)

// LoanRolePermissionsHandler defining the handler interface
type LoanRolePermissionsHandler interface {
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

type loanRolePermissionsHandler struct {
	iDao dao.LoanRolePermissionsDao
}

// NewLoanRolePermissionsHandler creating the handler interface
func NewLoanRolePermissionsHandler() LoanRolePermissionsHandler {
	return &loanRolePermissionsHandler{
		iDao: dao.NewLoanRolePermissionsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanRolePermissionsCache(database.GetCacheType()),
		),
	}
}

// Create a new loanRolePermissions
// @Summary Create a new loanRolePermissions
// @Description Creates a new loanRolePermissions entity using the provided data in the request body.
// @Tags loanRolePermissions
// @Accept json
// @Produce json
// @Param data body types.CreateLoanRolePermissionsRequest true "loanRolePermissions information"
// @Success 200 {object} types.CreateLoanRolePermissionsReply{}
// @Router /api/v1/loanRolePermissions [post]
// @Security BearerAuth
func (h *loanRolePermissionsHandler) Create(c *gin.Context) {
	form := &types.CreateLoanRolePermissionsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanRolePermissions := &model.LoanRolePermissions{}
	err = copier.Copy(loanRolePermissions, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanRolePermissions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanRolePermissions)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanRolePermissions.ID})
}

// DeleteByID delete a loanRolePermissions by id
// @Summary Delete a loanRolePermissions by id
// @Description Deletes a existing loanRolePermissions identified by the given id in the path.
// @Tags loanRolePermissions
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanRolePermissionsByIDReply{}
// @Router /api/v1/loanRolePermissions/{id} [delete]
// @Security BearerAuth
func (h *loanRolePermissionsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanRolePermissionsIDFromPath(c)
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

// UpdateByID update a loanRolePermissions by id
// @Summary Update a loanRolePermissions by id
// @Description Updates the specified loanRolePermissions by given id in the path, support partial update.
// @Tags loanRolePermissions
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanRolePermissionsByIDRequest true "loanRolePermissions information"
// @Success 200 {object} types.UpdateLoanRolePermissionsByIDReply{}
// @Router /api/v1/loanRolePermissions/{id} [put]
// @Security BearerAuth
func (h *loanRolePermissionsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanRolePermissionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanRolePermissionsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanRolePermissions := &model.LoanRolePermissions{}
	err = copier.Copy(loanRolePermissions, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanRolePermissions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanRolePermissions)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanRolePermissions by id
// @Summary Get a loanRolePermissions by id
// @Description Gets detailed information of a loanRolePermissions specified by the given id in the path.
// @Tags loanRolePermissions
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanRolePermissionsByIDReply{}
// @Router /api/v1/loanRolePermissions/{id} [get]
// @Security BearerAuth
func (h *loanRolePermissionsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanRolePermissionsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRolePermissions, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanRolePermissionsObjDetail{}
	err = copier.Copy(data, loanRolePermissions)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanRolePermissions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanRolePermissions": data})
}

// List get a paginated list of loanRolePermissionss by custom conditions
// @Summary Get a paginated list of loanRolePermissionss by custom conditions
// @Description Returns a paginated list of loanRolePermissions based on query filters, including page number and size.
// @Tags loanRolePermissions
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanRolePermissionssReply{}
// @Router /api/v1/loanRolePermissions/list [post]
// @Security BearerAuth
func (h *loanRolePermissionsHandler) List(c *gin.Context) {
	form := &types.ListLoanRolePermissionssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRolePermissionss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanRolePermissionss(loanRolePermissionss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanRolePermissions)
		return
	}

	response.Success(c, gin.H{
		"loanRolePermissionss": data,
		"total":                total,
	})
}

// DeleteByIDs batch delete loanRolePermissions by ids
// @Summary Batch delete loanRolePermissions by ids
// @Description Deletes multiple loanRolePermissions by a list of id
// @Tags loanRolePermissions
// @Param data body types.DeleteLoanRolePermissionssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanRolePermissionssByIDsReply{}
// @Router /api/v1/loanRolePermissions/delete/ids [post]
// @Security BearerAuth
func (h *loanRolePermissionsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanRolePermissionssByIDsRequest{}
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

// GetByCondition get a loanRolePermissions by custom condition
// @Summary Get a loanRolePermissions by custom condition
// @Description Returns a single loanRolePermissions that matches the specified filter conditions.
// @Tags loanRolePermissions
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanRolePermissionsByConditionReply{}
// @Router /api/v1/loanRolePermissions/condition [post]
// @Security BearerAuth
func (h *loanRolePermissionsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanRolePermissionsByConditionRequest{}
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
	loanRolePermissions, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanRolePermissionsObjDetail{}
	err = copier.Copy(data, loanRolePermissions)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanRolePermissions)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanRolePermissions": data})
}

// ListByIDs batch get loanRolePermissions by ids
// @Summary Batch get loanRolePermissions by ids
// @Description Returns a list of loanRolePermissions that match the list of id.
// @Tags loanRolePermissions
// @Param data body types.ListLoanRolePermissionssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanRolePermissionssByIDsReply{}
// @Router /api/v1/loanRolePermissions/list/ids [post]
// @Security BearerAuth
func (h *loanRolePermissionsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanRolePermissionssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRolePermissionsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanRolePermissionss := []*types.LoanRolePermissionsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanRolePermissionsMap[id]; ok {
			record, err := convertLoanRolePermissions(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanRolePermissions)
				return
			}
			loanRolePermissionss = append(loanRolePermissionss, record)
		}
	}

	response.Success(c, gin.H{
		"loanRolePermissionss": loanRolePermissionss,
	})
}

// ListByLastID get a paginated list of loanRolePermissionss by last id
// @Summary Get a paginated list of loanRolePermissionss by last id
// @Description Returns a paginated list of loanRolePermissionss starting after a given last id, useful for cursor-based pagination.
// @Tags loanRolePermissions
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanRolePermissionssReply{}
// @Router /api/v1/loanRolePermissions/list [get]
// @Security BearerAuth
func (h *loanRolePermissionsHandler) ListByLastID(c *gin.Context) {
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
	loanRolePermissionss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanRolePermissionss(loanRolePermissionss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanRolePermissions)
		return
	}

	response.Success(c, gin.H{
		"loanRolePermissionss": data,
	})
}

func getLoanRolePermissionsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanRolePermissions(loanRolePermissions *model.LoanRolePermissions) (*types.LoanRolePermissionsObjDetail, error) {
	data := &types.LoanRolePermissionsObjDetail{}
	err := copier.Copy(data, loanRolePermissions)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanRolePermissionss(fromValues []*model.LoanRolePermissions) ([]*types.LoanRolePermissionsObjDetail, error) {
	toValues := []*types.LoanRolePermissionsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanRolePermissions(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
