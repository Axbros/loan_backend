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

var _ LoanUserCallRecordsHandler = (*loanUserCallRecordsHandler)(nil)

// LoanUserCallRecordsHandler defining the handler interface
type LoanUserCallRecordsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type loanUserCallRecordsHandler struct {
	iDao dao.LoanUserCallRecordsDao
}

// NewLoanUserCallRecordsHandler creating the handler interface
func NewLoanUserCallRecordsHandler() LoanUserCallRecordsHandler {
	return &loanUserCallRecordsHandler{
		iDao: dao.NewLoanUserCallRecordsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanUserCallRecordsCache(database.GetCacheType()),
		),
	}
}

// Create a new loanUserCallRecords
// @Summary Create a new loanUserCallRecords
// @Description Creates a new loanUserCallRecords entity using the provided data in the request body.
// @Tags loanUserCallRecords
// @Accept json
// @Produce json
// @Param data body types.CreateLoanUserCallRecordsRequest true "loanUserCallRecords information"
// @Success 200 {object} types.CreateLoanUserCallRecordsReply{}
// @Router /api/v1/loanUserCallRecords [post]
// @Security BearerAuth
func (h *loanUserCallRecordsHandler) Create(c *gin.Context) {
	form := &types.CreateLoanUserCallRecordsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanUserCallRecords := &model.LoanUserCallRecords{}
	err = copier.Copy(loanUserCallRecords, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanUserCallRecords)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanUserCallRecords)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanUserCallRecords.ID})
}

// DeleteByID delete a loanUserCallRecords by id
// @Summary Delete a loanUserCallRecords by id
// @Description Deletes a existing loanUserCallRecords identified by the given id in the path.
// @Tags loanUserCallRecords
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanUserCallRecordsByIDReply{}
// @Router /api/v1/loanUserCallRecords/{id} [delete]
// @Security BearerAuth
func (h *loanUserCallRecordsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanUserCallRecordsIDFromPath(c)
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

// UpdateByID update a loanUserCallRecords by id
// @Summary Update a loanUserCallRecords by id
// @Description Updates the specified loanUserCallRecords by given id in the path, support partial update.
// @Tags loanUserCallRecords
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanUserCallRecordsByIDRequest true "loanUserCallRecords information"
// @Success 200 {object} types.UpdateLoanUserCallRecordsByIDReply{}
// @Router /api/v1/loanUserCallRecords/{id} [put]
// @Security BearerAuth
func (h *loanUserCallRecordsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanUserCallRecordsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanUserCallRecordsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanUserCallRecords := &model.LoanUserCallRecords{}
	err = copier.Copy(loanUserCallRecords, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanUserCallRecords)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanUserCallRecords)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanUserCallRecords by id
// @Summary Get a loanUserCallRecords by id
// @Description Gets detailed information of a loanUserCallRecords specified by the given id in the path.
// @Tags loanUserCallRecords
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanUserCallRecordsByIDReply{}
// @Router /api/v1/loanUserCallRecords/{id} [get]
// @Security BearerAuth
func (h *loanUserCallRecordsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanUserCallRecordsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUserCallRecords, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanUserCallRecordsObjDetail{}
	err = copier.Copy(data, loanUserCallRecords)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanUserCallRecords)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanUserCallRecords": data})
}

// List get a paginated list of loanUserCallRecordss by custom conditions
// @Summary Get a paginated list of loanUserCallRecordss by custom conditions
// @Description Returns a paginated list of loanUserCallRecords based on query filters, including page number and size.
// @Tags loanUserCallRecords
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanUserCallRecordssReply{}
// @Router /api/v1/loanUserCallRecords/list [post]
// @Security BearerAuth
func (h *loanUserCallRecordsHandler) List(c *gin.Context) {
	form := &types.ListLoanUserCallRecordssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUserCallRecordss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanUserCallRecordss(loanUserCallRecordss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanUserCallRecords)
		return
	}

	response.Success(c, gin.H{
		"records": data,
		"total":   total,
	})
}

func getLoanUserCallRecordsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanUserCallRecords(loanUserCallRecords *model.LoanUserCallRecords) (*types.LoanUserCallRecordsObjDetail, error) {
	data := &types.LoanUserCallRecordsObjDetail{}
	err := copier.Copy(data, loanUserCallRecords)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanUserCallRecordss(fromValues []*model.LoanUserCallRecords) ([]*types.LoanUserCallRecordsObjDetail, error) {
	toValues := []*types.LoanUserCallRecordsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanUserCallRecords(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
