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

var _ LoanRolesHandler = (*loanRolesHandler)(nil)

// LoanRolesHandler defining the handler interface
type LoanRolesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
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
		"records": data,
		"total":   total,
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
