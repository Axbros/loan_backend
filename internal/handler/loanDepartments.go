package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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

var _ LoanDepartmentsHandler = (*loanDepartmentsHandler)(nil)

// LoanDepartmentsHandler defining the handler interface
type LoanDepartmentsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type loanDepartmentsHandler struct {
	iDao              dao.LoanDepartmentsDao
	roleDepartmentDao dao.LoanDepartmentRolesDao
}

// NewLoanDepartmentsHandler creating the handler interface
func NewLoanDepartmentsHandler() LoanDepartmentsHandler {
	return &loanDepartmentsHandler{
		iDao: dao.NewLoanDepartmentsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanDepartmentsCache(database.GetCacheType()),
		),
		roleDepartmentDao: dao.NewLoanDepartmentRolesDao(
			database.GetDB(),
			cache.NewLoanDepartmentRolesCache(database.GetCacheType()),
		),
	}
}

// Create a new loanDepartments
// @Summary Create a new loanDepartments
// @Description Creates a new loanDepartments entity using the provided data in the request body.
// @Tags loanDepartments
// @Accept json
// @Produce json
// @Param data body types.CreateLoanDepartmentsRequest true "loanDepartments information"
// @Success 200 {object} types.CreateLoanDepartmentsReply{}
// @Router /api/v1/loanDepartments [post]
// @Security BearerAuth
func (h *loanDepartmentsHandler) Create(c *gin.Context) {
	form := &types.CreateLoanDepartmentsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanDepartments := &model.LoanDepartments{}
	err = copier.Copy(loanDepartments, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanDepartments)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanDepartments)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanDepartments.ID})
}

// DeleteByID delete a loanDepartments by id
// @Summary Delete a loanDepartments by id
// @Description Deletes a existing loanDepartments identified by the given id in the path.
// @Tags loanDepartments
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanDepartmentsByIDReply{}
// @Router /api/v1/loanDepartments/{id} [delete]
// @Security BearerAuth
func (h *loanDepartmentsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanDepartmentsIDFromPath(c)
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

// UpdateByID update a loanDepartments by id
// @Summary Update a loanDepartments by id
// @Description Updates the specified loanDepartments by given id in the path, support partial update.
// @Tags loanDepartments
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanDepartmentsByIDRequest true "loanDepartments information"
// @Success 200 {object} types.UpdateLoanDepartmentsByIDReply{}
// @Router /api/v1/loanDepartments/{id} [put]
// @Security BearerAuth
func (h *loanDepartmentsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanDepartmentsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanDepartmentsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanDepartments := &model.LoanDepartments{}
	err = copier.Copy(loanDepartments, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanDepartments)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanDepartments)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	record, err := h.roleDepartmentDao.GetByDepartmentID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在就新增
			newRecord := &model.LoanDepartmentRoles{
				DepartmentID: id,
				RoleID:       form.RoleID,
			}
			if err := h.roleDepartmentDao.Create(ctx, newRecord); err != nil {
				logger.Error("Create LoanDepartmentRoles error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
				response.Error(c, ecode.ErrCreateLoanDepartmentRoles)
				return
			}
		} else {
			logger.Error("GetByDepartmentID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.ErrGetByConditionLoanDepartmentRoles)
			return
		}
	} else {
		// 存在就更新
		record.RoleID = form.RoleID
		if err := h.roleDepartmentDao.UpdateByID(ctx, record); err != nil {
			logger.Error("UpdateByID LoanDepartmentRoles error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.ErrUpdateByIDLoanDepartmentRoles)
			return
		}
	}
	response.Success(c)
}

// GetByID get a loanDepartments by id
// @Summary Get a loanDepartments by id
// @Description Gets detailed information of a loanDepartments specified by the given id in the path.
// @Tags loanDepartments
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanDepartmentsByIDReply{}
// @Router /api/v1/loanDepartments/{id} [get]
// @Security BearerAuth
func (h *loanDepartmentsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanDepartmentsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanDepartments, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanDepartmentsObjDetail{}
	err = copier.Copy(data, loanDepartments)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanDepartments)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanDepartments": data})
}

// List get a paginated list of loanDepartmentss by custom conditions
// @Summary Get a paginated list of loanDepartmentss by custom conditions
// @Description Returns a paginated list of loanDepartments based on query filters, including page number and size.
// @Tags loanDepartments
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanDepartmentssReply{}
// @Router /api/v1/loanDepartments/list [post]
// @Security BearerAuth
func (h *loanDepartmentsHandler) List(c *gin.Context) {
	form := &types.ListLoanDepartmentssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanDepartmentss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{
		"records": loanDepartmentss,
		"total":   total,
	})
}

func getLoanDepartmentsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanDepartments(loanDepartments *model.LoanDepartments) (*types.LoanDepartmentsObjDetail, error) {
	data := &types.LoanDepartmentsObjDetail{}
	err := copier.Copy(data, loanDepartments)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanDepartmentss(fromValues []*model.LoanDepartments) ([]*types.LoanDepartmentsObjDetail, error) {
	toValues := []*types.LoanDepartmentsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanDepartments(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
