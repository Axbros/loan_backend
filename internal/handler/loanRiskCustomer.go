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

var _ LoanRiskCustomerHandler = (*loanRiskCustomerHandler)(nil)

// LoanRiskCustomerHandler defining the handler interface
type LoanRiskCustomerHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type loanRiskCustomerHandler struct {
	iDao dao.LoanRiskCustomerDao
}

// NewLoanRiskCustomerHandler creating the handler interface
func NewLoanRiskCustomerHandler() LoanRiskCustomerHandler {
	return &loanRiskCustomerHandler{
		iDao: dao.NewLoanRiskCustomerDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanRiskCustomerCache(database.GetCacheType()),
		),
	}
}

// Create a new loanRiskCustomer
// @Summary Create a new loanRiskCustomer
// @Description Creates a new loanRiskCustomer entity using the provided data in the request body.
// @Tags loanRiskCustomer
// @Accept json
// @Produce json
// @Param data body types.CreateLoanRiskCustomerRequest true "loanRiskCustomer information"
// @Success 200 {object} types.CreateLoanRiskCustomerReply{}
// @Router /api/v1/loanRiskCustomer [post]
// @Security BearerAuth
func (h *loanRiskCustomerHandler) Create(c *gin.Context) {
	form := &types.CreateLoanRiskCustomerRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanRiskCustomer := &model.LoanRiskCustomer{}
	err = copier.Copy(loanRiskCustomer, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanRiskCustomer)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanRiskCustomer)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanRiskCustomer.ID})
}

// DeleteByID delete a loanRiskCustomer by id
// @Summary Delete a loanRiskCustomer by id
// @Description Deletes a existing loanRiskCustomer identified by the given id in the path.
// @Tags loanRiskCustomer
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanRiskCustomerByIDReply{}
// @Router /api/v1/loanRiskCustomer/{id} [delete]
// @Security BearerAuth
func (h *loanRiskCustomerHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanRiskCustomerIDFromPath(c)
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

// UpdateByID update a loanRiskCustomer by id
// @Summary Update a loanRiskCustomer by id
// @Description Updates the specified loanRiskCustomer by given id in the path, support partial update.
// @Tags loanRiskCustomer
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanRiskCustomerByIDRequest true "loanRiskCustomer information"
// @Success 200 {object} types.UpdateLoanRiskCustomerByIDReply{}
// @Router /api/v1/loanRiskCustomer/{id} [put]
// @Security BearerAuth
func (h *loanRiskCustomerHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanRiskCustomerIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanRiskCustomerByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanRiskCustomer := &model.LoanRiskCustomer{}
	err = copier.Copy(loanRiskCustomer, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanRiskCustomer)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanRiskCustomer)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanRiskCustomer by id
// @Summary Get a loanRiskCustomer by id
// @Description Gets detailed information of a loanRiskCustomer specified by the given id in the path.
// @Tags loanRiskCustomer
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanRiskCustomerByIDReply{}
// @Router /api/v1/loanRiskCustomer/{id} [get]
// @Security BearerAuth
func (h *loanRiskCustomerHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanRiskCustomerIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRiskCustomer, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanRiskCustomerObjDetail{}
	err = copier.Copy(data, loanRiskCustomer)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanRiskCustomer)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanRiskCustomer": data})
}

// List get a paginated list of loanRiskCustomers by custom conditions
// @Summary Get a paginated list of loanRiskCustomers by custom conditions
// @Description Returns a paginated list of loanRiskCustomer based on query filters, including page number and size.
// @Tags loanRiskCustomer
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanRiskCustomersReply{}
// @Router /api/v1/loanRiskCustomer/list [post]
// @Security BearerAuth
func (h *loanRiskCustomerHandler) List(c *gin.Context) {
	form := &types.ListLoanRiskCustomersRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanRiskCustomers, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanRiskCustomers(loanRiskCustomers)
	if err != nil {
		response.Error(c, ecode.ErrListLoanRiskCustomer)
		return
	}

	response.Success(c, gin.H{
		"loanRiskCustomers": data,
		"total":        total,
	})
}

func getLoanRiskCustomerIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanRiskCustomer(loanRiskCustomer *model.LoanRiskCustomer) (*types.LoanRiskCustomerObjDetail, error) {
	data := &types.LoanRiskCustomerObjDetail{}
	err := copier.Copy(data, loanRiskCustomer)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanRiskCustomers(fromValues []*model.LoanRiskCustomer) ([]*types.LoanRiskCustomerObjDetail, error) {
	toValues := []*types.LoanRiskCustomerObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanRiskCustomer(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
