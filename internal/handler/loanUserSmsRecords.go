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

var _ LoanUserSmsRecordsHandler = (*loanUserSmsRecordsHandler)(nil)

// LoanUserSmsRecordsHandler defining the handler interface
type LoanUserSmsRecordsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type loanUserSmsRecordsHandler struct {
	iDao dao.LoanUserSmsRecordsDao
}

// NewLoanUserSmsRecordsHandler creating the handler interface
func NewLoanUserSmsRecordsHandler() LoanUserSmsRecordsHandler {
	return &loanUserSmsRecordsHandler{
		iDao: dao.NewLoanUserSmsRecordsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanUserSmsRecordsCache(database.GetCacheType()),
		),
	}
}

// Create a new loanUserSmsRecords
// @Summary Create a new loanUserSmsRecords
// @Description Creates a new loanUserSmsRecords entity using the provided data in the request body.
// @Tags loanUserSmsRecords
// @Accept json
// @Produce json
// @Param data body types.CreateLoanUserSmsRecordsRequest true "loanUserSmsRecords information"
// @Success 200 {object} types.CreateLoanUserSmsRecordsReply{}
// @Router /api/v1/loanUserSmsRecords [post]
// @Security BearerAuth
func (h *loanUserSmsRecordsHandler) Create(c *gin.Context) {
	form := &types.CreateLoanUserSmsRecordsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanUserSmsRecords := &model.LoanUserSmsRecords{}
	err = copier.Copy(loanUserSmsRecords, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanUserSmsRecords)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanUserSmsRecords)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanUserSmsRecords.ID})
}

// DeleteByID delete a loanUserSmsRecords by id
// @Summary Delete a loanUserSmsRecords by id
// @Description Deletes a existing loanUserSmsRecords identified by the given id in the path.
// @Tags loanUserSmsRecords
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanUserSmsRecordsByIDReply{}
// @Router /api/v1/loanUserSmsRecords/{id} [delete]
// @Security BearerAuth
func (h *loanUserSmsRecordsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanUserSmsRecordsIDFromPath(c)
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

// UpdateByID update a loanUserSmsRecords by id
// @Summary Update a loanUserSmsRecords by id
// @Description Updates the specified loanUserSmsRecords by given id in the path, support partial update.
// @Tags loanUserSmsRecords
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanUserSmsRecordsByIDRequest true "loanUserSmsRecords information"
// @Success 200 {object} types.UpdateLoanUserSmsRecordsByIDReply{}
// @Router /api/v1/loanUserSmsRecords/{id} [put]
// @Security BearerAuth
func (h *loanUserSmsRecordsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanUserSmsRecordsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanUserSmsRecordsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanUserSmsRecords := &model.LoanUserSmsRecords{}
	err = copier.Copy(loanUserSmsRecords, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanUserSmsRecords)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanUserSmsRecords)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanUserSmsRecords by id
// @Summary Get a loanUserSmsRecords by id
// @Description Gets detailed information of a loanUserSmsRecords specified by the given id in the path.
// @Tags loanUserSmsRecords
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanUserSmsRecordsByIDReply{}
// @Router /api/v1/loanUserSmsRecords/{id} [get]
// @Security BearerAuth
func (h *loanUserSmsRecordsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanUserSmsRecordsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUserSmsRecords, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanUserSmsRecordsObjDetail{}
	err = copier.Copy(data, loanUserSmsRecords)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanUserSmsRecords)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanUserSmsRecords": data})
}

// List get a paginated list of loanUserSmsRecordss by custom conditions
// @Summary Get a paginated list of loanUserSmsRecordss by custom conditions
// @Description Returns a paginated list of loanUserSmsRecords based on query filters, including page number and size.
// @Tags loanUserSmsRecords
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanUserSmsRecordssReply{}
// @Router /api/v1/loanUserSmsRecords/list [post]
// @Security BearerAuth
func (h *loanUserSmsRecordsHandler) List(c *gin.Context) {
	form := &types.ListLoanUserSmsRecordssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUserSmsRecordss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanUserSmsRecordss(loanUserSmsRecordss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanUserSmsRecords)
		return
	}

	response.Success(c, gin.H{
		"records": data,
		"total":   total,
	})
}

func getLoanUserSmsRecordsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanUserSmsRecords(loanUserSmsRecords *model.LoanUserSmsRecords) (*types.LoanUserSmsRecordsObjDetail, error) {
	data := &types.LoanUserSmsRecordsObjDetail{}
	err := copier.Copy(data, loanUserSmsRecords)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanUserSmsRecordss(fromValues []*model.LoanUserSmsRecords) ([]*types.LoanUserSmsRecordsObjDetail, error) {
	toValues := []*types.LoanUserSmsRecordsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanUserSmsRecords(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
