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

var _ LoanRolesHandler = (*loanRolesHandler)(nil)

// LoanRolesHandler defining the handler interface
type LoanRolesHandler interface {
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

type loanRolesHandler struct {
	iDao dao.LoanRolesDao
}

// NewLoanRolesHandler creating the handler interface
func NewLoanRolesHandler() LoanRolesHandler {
	return &loanRolesHandler{
		iDao: dao.NewLoanRolesDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanRolesCache(database.GetCacheType()),
		),
	}
}

// Create a new loanRoles
// @Summary Create a new loanRoles
// @Description Creates a new loanRoles entity using the provided data in the request body.
// @Tags loanRoles
// @Accept json
// @Produce json
// @Param data body types.CreateLoanRolesRequest true "loanRoles information"
// @Success 200 {object} types.CreateLoanRolesReply{}
// @Router /api/v1/loanRoles [post]
// @Security BearerAuth
func (h *loanRolesHandler) Create(c *gin.Context) {
	form := &types.CreateLoanRolesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanRoles := &model.LoanRoles{}
	err = copier.Copy(loanRoles, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanRoles)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanRoles.ID})
}

// DeleteByID delete a loanRoles by id
// @Summary Delete a loanRoles by id
// @Description Deletes a existing loanRoles identified by the given id in the path.
// @Tags loanRoles
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanRolesByIDReply{}
// @Router /api/v1/loanRoles/{id} [delete]
// @Security BearerAuth
func (h *loanRolesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanRolesIDFromPath(c)
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

// UpdateByID update a loanRoles by id
// @Summary Update a loanRoles by id
// @Description Updates the specified loanRoles by given id in the path, support partial update.
// @Tags loanRoles
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanRolesByIDRequest true "loanRoles information"
// @Success 200 {object} types.UpdateLoanRolesByIDReply{}
// @Router /api/v1/loanRoles/{id} [put]
// @Security BearerAuth
func (h *loanRolesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanRolesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanRolesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanRoles := &model.LoanRoles{}
	err = copier.Copy(loanRoles, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanRoles)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanRoles by id
// @Summary Get a loanRoles by id
// @Description Gets detailed information of a loanRoles specified by the given id in the path.
// @Tags loanRoles
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanRolesByIDReply{}
// @Router /api/v1/loanRoles/{id} [get]
// @Security BearerAuth
func (h *loanRolesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanRolesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRoles, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanRolesObjDetail{}
	err = copier.Copy(data, loanRoles)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanRoles": data})
}

// List get a paginated list of loanRoless by custom conditions
// @Summary Get a paginated list of loanRoless by custom conditions
// @Description Returns a paginated list of loanRoles based on query filters, including page number and size.
// @Tags loanRoles
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanRolessReply{}
// @Router /api/v1/loanRoles/list [post]
// @Security BearerAuth
func (h *loanRolesHandler) List(c *gin.Context) {
	form := &types.ListLoanRolessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRoless, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanRoless(loanRoless)
	if err != nil {
		response.Error(c, ecode.ErrListLoanRoles)
		return
	}

	response.Success(c, gin.H{
		"loanRoless": data,
		"total":      total,
	})
}

// DeleteByIDs batch delete loanRoles by ids
// @Summary Batch delete loanRoles by ids
// @Description Deletes multiple loanRoles by a list of id
// @Tags loanRoles
// @Param data body types.DeleteLoanRolessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanRolessByIDsReply{}
// @Router /api/v1/loanRoles/delete/ids [post]
// @Security BearerAuth
func (h *loanRolesHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanRolessByIDsRequest{}
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

// GetByCondition get a loanRoles by custom condition
// @Summary Get a loanRoles by custom condition
// @Description Returns a single loanRoles that matches the specified filter conditions.
// @Tags loanRoles
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanRolesByConditionReply{}
// @Router /api/v1/loanRoles/condition [post]
// @Security BearerAuth
func (h *loanRolesHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanRolesByConditionRequest{}
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
	loanRoles, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanRolesObjDetail{}
	err = copier.Copy(data, loanRoles)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanRoles": data})
}

// ListByIDs batch get loanRoles by ids
// @Summary Batch get loanRoles by ids
// @Description Returns a list of loanRoles that match the list of id.
// @Tags loanRoles
// @Param data body types.ListLoanRolessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanRolessByIDsReply{}
// @Router /api/v1/loanRoles/list/ids [post]
// @Security BearerAuth
func (h *loanRolesHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanRolessByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRolesMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanRoless := []*types.LoanRolesObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanRolesMap[id]; ok {
			record, err := convertLoanRoles(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanRoles)
				return
			}
			loanRoless = append(loanRoless, record)
		}
	}

	response.Success(c, gin.H{
		"loanRoless": loanRoless,
	})
}

// ListByLastID get a paginated list of loanRoless by last id
// @Summary Get a paginated list of loanRoless by last id
// @Description Returns a paginated list of loanRoless starting after a given last id, useful for cursor-based pagination.
// @Tags loanRoles
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanRolessReply{}
// @Router /api/v1/loanRoles/list [get]
// @Security BearerAuth
func (h *loanRolesHandler) ListByLastID(c *gin.Context) {
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
	loanRoless, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanRoless(loanRoless)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanRoles)
		return
	}

	response.Success(c, gin.H{
		"loanRoless": data,
	})
}

func getLoanRolesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanRoles(loanRoles *model.LoanRoles) (*types.LoanRolesObjDetail, error) {
	data := &types.LoanRolesObjDetail{}
	err := copier.Copy(data, loanRoles)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanRoless(fromValues []*model.LoanRoles) ([]*types.LoanRolesObjDetail, error) {
	toValues := []*types.LoanRolesObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanRoles(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
