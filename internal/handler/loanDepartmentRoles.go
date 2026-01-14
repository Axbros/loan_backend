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

var _ LoanDepartmentRolesHandler = (*loanDepartmentRolesHandler)(nil)

// LoanDepartmentRolesHandler defining the handler interface
type LoanDepartmentRolesHandler interface {
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

type loanDepartmentRolesHandler struct {
	iDao dao.LoanDepartmentRolesDao
}

// NewLoanDepartmentRolesHandler creating the handler interface
func NewLoanDepartmentRolesHandler() LoanDepartmentRolesHandler {
	return &loanDepartmentRolesHandler{
		iDao: dao.NewLoanDepartmentRolesDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanDepartmentRolesCache(database.GetCacheType()),
		),
	}
}

// Create a new loanDepartmentRoles
// @Summary Create a new loanDepartmentRoles
// @Description Creates a new loanDepartmentRoles entity using the provided data in the request body.
// @Tags loanDepartmentRoles
// @Accept json
// @Produce json
// @Param data body types.CreateLoanDepartmentRolesRequest true "loanDepartmentRoles information"
// @Success 200 {object} types.CreateLoanDepartmentRolesReply{}
// @Router /api/v1/loanDepartmentRoles [post]
// @Security BearerAuth
func (h *loanDepartmentRolesHandler) Create(c *gin.Context) {
	form := &types.CreateLoanDepartmentRolesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanDepartmentRoles := &model.LoanDepartmentRoles{}
	err = copier.Copy(loanDepartmentRoles, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanDepartmentRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanDepartmentRoles)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanDepartmentRoles.ID})
}

// DeleteByID delete a loanDepartmentRoles by id
// @Summary Delete a loanDepartmentRoles by id
// @Description Deletes a existing loanDepartmentRoles identified by the given id in the path.
// @Tags loanDepartmentRoles
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanDepartmentRolesByIDReply{}
// @Router /api/v1/loanDepartmentRoles/{id} [delete]
// @Security BearerAuth
func (h *loanDepartmentRolesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanDepartmentRolesIDFromPath(c)
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

// UpdateByID update a loanDepartmentRoles by id
// @Summary Update a loanDepartmentRoles by id
// @Description Updates the specified loanDepartmentRoles by given id in the path, support partial update.
// @Tags loanDepartmentRoles
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanDepartmentRolesByIDRequest true "loanDepartmentRoles information"
// @Success 200 {object} types.UpdateLoanDepartmentRolesByIDReply{}
// @Router /api/v1/loanDepartmentRoles/{id} [put]
// @Security BearerAuth
func (h *loanDepartmentRolesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanDepartmentRolesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanDepartmentRolesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanDepartmentRoles := &model.LoanDepartmentRoles{}
	err = copier.Copy(loanDepartmentRoles, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanDepartmentRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanDepartmentRoles)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanDepartmentRoles by id
// @Summary Get a loanDepartmentRoles by id
// @Description Gets detailed information of a loanDepartmentRoles specified by the given id in the path.
// @Tags loanDepartmentRoles
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanDepartmentRolesByIDReply{}
// @Router /api/v1/loanDepartmentRoles/{id} [get]
// @Security BearerAuth
func (h *loanDepartmentRolesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanDepartmentRolesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanDepartmentRoles, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanDepartmentRolesObjDetail{}
	err = copier.Copy(data, loanDepartmentRoles)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanDepartmentRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanDepartmentRoles": data})
}

// List get a paginated list of loanDepartmentRoless by custom conditions
// @Summary Get a paginated list of loanDepartmentRoless by custom conditions
// @Description Returns a paginated list of loanDepartmentRoles based on query filters, including page number and size.
// @Tags loanDepartmentRoles
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanDepartmentRolessReply{}
// @Router /api/v1/loanDepartmentRoles/list [post]
// @Security BearerAuth
func (h *loanDepartmentRolesHandler) List(c *gin.Context) {
	form := &types.ListLoanDepartmentRolessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanDepartmentRoless, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanDepartmentRoless(loanDepartmentRoless)
	if err != nil {
		response.Error(c, ecode.ErrListLoanDepartmentRoles)
		return
	}

	response.Success(c, gin.H{
		"loanDepartmentRoless": data,
		"total":                total,
	})
}

// DeleteByIDs batch delete loanDepartmentRoles by ids
// @Summary Batch delete loanDepartmentRoles by ids
// @Description Deletes multiple loanDepartmentRoles by a list of id
// @Tags loanDepartmentRoles
// @Param data body types.DeleteLoanDepartmentRolessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanDepartmentRolessByIDsReply{}
// @Router /api/v1/loanDepartmentRoles/delete/ids [post]
// @Security BearerAuth
func (h *loanDepartmentRolesHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanDepartmentRolessByIDsRequest{}
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

// GetByCondition get a loanDepartmentRoles by custom condition
// @Summary Get a loanDepartmentRoles by custom condition
// @Description Returns a single loanDepartmentRoles that matches the specified filter conditions.
// @Tags loanDepartmentRoles
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanDepartmentRolesByConditionReply{}
// @Router /api/v1/loanDepartmentRoles/condition [post]
// @Security BearerAuth
func (h *loanDepartmentRolesHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanDepartmentRolesByConditionRequest{}
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
	loanDepartmentRoles, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanDepartmentRolesObjDetail{}
	err = copier.Copy(data, loanDepartmentRoles)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanDepartmentRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanDepartmentRoles": data})
}

// ListByIDs batch get loanDepartmentRoles by ids
// @Summary Batch get loanDepartmentRoles by ids
// @Description Returns a list of loanDepartmentRoles that match the list of id.
// @Tags loanDepartmentRoles
// @Param data body types.ListLoanDepartmentRolessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanDepartmentRolessByIDsReply{}
// @Router /api/v1/loanDepartmentRoles/list/ids [post]
// @Security BearerAuth
func (h *loanDepartmentRolesHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanDepartmentRolessByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanDepartmentRolesMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanDepartmentRoless := []*types.LoanDepartmentRolesObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanDepartmentRolesMap[id]; ok {
			record, err := convertLoanDepartmentRoles(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanDepartmentRoles)
				return
			}
			loanDepartmentRoless = append(loanDepartmentRoless, record)
		}
	}

	response.Success(c, gin.H{
		"loanDepartmentRoless": loanDepartmentRoless,
	})
}

// ListByLastID get a paginated list of loanDepartmentRoless by last id
// @Summary Get a paginated list of loanDepartmentRoless by last id
// @Description Returns a paginated list of loanDepartmentRoless starting after a given last id, useful for cursor-based pagination.
// @Tags loanDepartmentRoles
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanDepartmentRolessReply{}
// @Router /api/v1/loanDepartmentRoles/list [get]
// @Security BearerAuth
func (h *loanDepartmentRolesHandler) ListByLastID(c *gin.Context) {
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
	loanDepartmentRoless, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanDepartmentRoless(loanDepartmentRoless)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanDepartmentRoles)
		return
	}

	response.Success(c, gin.H{
		"loanDepartmentRoless": data,
	})
}

func getLoanDepartmentRolesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanDepartmentRoles(loanDepartmentRoles *model.LoanDepartmentRoles) (*types.LoanDepartmentRolesObjDetail, error) {
	data := &types.LoanDepartmentRolesObjDetail{}
	err := copier.Copy(data, loanDepartmentRoles)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanDepartmentRoless(fromValues []*model.LoanDepartmentRoles) ([]*types.LoanDepartmentRolesObjDetail, error) {
	toValues := []*types.LoanDepartmentRolesObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanDepartmentRoles(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
