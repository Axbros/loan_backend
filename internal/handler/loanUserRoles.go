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

var _ LoanUserRolesHandler = (*loanUserRolesHandler)(nil)

// LoanUserRolesHandler defining the handler interface
type LoanUserRolesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type loanUserRolesHandler struct {
	iDao dao.LoanUserRolesDao
}

// NewLoanUserRolesHandler creating the handler interface
func NewLoanUserRolesHandler() LoanUserRolesHandler {
	return &loanUserRolesHandler{
		iDao: dao.NewLoanUserRolesDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanUserRolesCache(database.GetCacheType()),
		),
	}
}

// Create a new loanUserRoles
// @Summary Create a new loanUserRoles
// @Description Creates a new loanUserRoles entity using the provided data in the request body.
// @Tags loanUserRoles
// @Accept json
// @Produce json
// @Param data body types.CreateLoanUserRolesRequest true "loanUserRoles information"
// @Success 200 {object} types.CreateLoanUserRolesReply{}
// @Router /api/v1/loanUserRoles [post]
// @Security BearerAuth
func (h *loanUserRolesHandler) Create(c *gin.Context) {
	form := &types.CreateLoanUserRolesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanUserRoles := &model.LoanUserRoles{}
	err = copier.Copy(loanUserRoles, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanUserRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanUserRoles)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanUserRoles.ID})
}

// DeleteByID delete a loanUserRoles by id
// @Summary Delete a loanUserRoles by id
// @Description Deletes a existing loanUserRoles identified by the given id in the path.
// @Tags loanUserRoles
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanUserRolesByIDReply{}
// @Router /api/v1/loanUserRoles/{id} [delete]
// @Security BearerAuth
func (h *loanUserRolesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanUserRolesIDFromPath(c)
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

// UpdateByID update a loanUserRoles by id
// @Summary Update a loanUserRoles by id
// @Description Updates the specified loanUserRoles by given id in the path, support partial update.
// @Tags loanUserRoles
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanUserRolesByIDRequest true "loanUserRoles information"
// @Success 200 {object} types.UpdateLoanUserRolesByIDReply{}
// @Router /api/v1/loanUserRoles/{id} [put]
// @Security BearerAuth
func (h *loanUserRolesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanUserRolesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanUserRolesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanUserRoles := &model.LoanUserRoles{}
	err = copier.Copy(loanUserRoles, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanUserRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanUserRoles)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanUserRoles by id
// @Summary Get a loanUserRoles by id
// @Description Gets detailed information of a loanUserRoles specified by the given id in the path.
// @Tags loanUserRoles
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanUserRolesByIDReply{}
// @Router /api/v1/loanUserRoles/{id} [get]
// @Security BearerAuth
func (h *loanUserRolesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanUserRolesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUserRoles, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanUserRolesObjDetail{}
	err = copier.Copy(data, loanUserRoles)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanUserRoles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanUserRoles": data})
}

// List get a paginated list of loanUserRoless by custom conditions
// @Summary Get a paginated list of loanUserRoless by custom conditions
// @Description Returns a paginated list of loanUserRoles based on query filters, including page number and size.
// @Tags loanUserRoles
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanUserRolessReply{}
// @Router /api/v1/loanUserRoles/list [post]
// @Security BearerAuth
func (h *loanUserRolesHandler) List(c *gin.Context) {
	form := &types.ListLoanUserRolessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUserRoless, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanUserRoless(loanUserRoless)
	if err != nil {
		response.Error(c, ecode.ErrListLoanUserRoles)
		return
	}

	response.Success(c, gin.H{
		"loanUserRoless": data,
		"total":          total,
	})
}

func getLoanUserRolesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanUserRoles(loanUserRoles *model.LoanUserRoles) (*types.LoanUserRolesObjDetail, error) {
	data := &types.LoanUserRolesObjDetail{}
	err := copier.Copy(data, loanUserRoles)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanUserRoless(fromValues []*model.LoanUserRoles) ([]*types.LoanUserRolesObjDetail, error) {
	toValues := []*types.LoanUserRolesObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanUserRoles(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
