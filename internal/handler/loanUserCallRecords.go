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

var _ LoanUserCallRecordsHandler = (*loanUserCallRecordsHandler)(nil)

// LoanUserCallRecordsHandler defining the handler interface
type LoanUserCallRecordsHandler interface {
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
		"loanUserCallRecordss": data,
		"total":                total,
	})
}

// DeleteByIDs batch delete loanUserCallRecords by ids
// @Summary Batch delete loanUserCallRecords by ids
// @Description Deletes multiple loanUserCallRecords by a list of id
// @Tags loanUserCallRecords
// @Param data body types.DeleteLoanUserCallRecordssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanUserCallRecordssByIDsReply{}
// @Router /api/v1/loanUserCallRecords/delete/ids [post]
// @Security BearerAuth
func (h *loanUserCallRecordsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanUserCallRecordssByIDsRequest{}
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

// GetByCondition get a loanUserCallRecords by custom condition
// @Summary Get a loanUserCallRecords by custom condition
// @Description Returns a single loanUserCallRecords that matches the specified filter conditions.
// @Tags loanUserCallRecords
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanUserCallRecordsByConditionReply{}
// @Router /api/v1/loanUserCallRecords/condition [post]
// @Security BearerAuth
func (h *loanUserCallRecordsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanUserCallRecordsByConditionRequest{}
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
	loanUserCallRecords, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanUserCallRecordsObjDetail{}
	err = copier.Copy(data, loanUserCallRecords)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanUserCallRecords)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanUserCallRecords": data})
}

// ListByIDs batch get loanUserCallRecords by ids
// @Summary Batch get loanUserCallRecords by ids
// @Description Returns a list of loanUserCallRecords that match the list of id.
// @Tags loanUserCallRecords
// @Param data body types.ListLoanUserCallRecordssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanUserCallRecordssByIDsReply{}
// @Router /api/v1/loanUserCallRecords/list/ids [post]
// @Security BearerAuth
func (h *loanUserCallRecordsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanUserCallRecordssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUserCallRecordsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanUserCallRecordss := []*types.LoanUserCallRecordsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanUserCallRecordsMap[id]; ok {
			record, err := convertLoanUserCallRecords(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanUserCallRecords)
				return
			}
			loanUserCallRecordss = append(loanUserCallRecordss, record)
		}
	}

	response.Success(c, gin.H{
		"loanUserCallRecordss": loanUserCallRecordss,
	})
}

// ListByLastID get a paginated list of loanUserCallRecordss by last id
// @Summary Get a paginated list of loanUserCallRecordss by last id
// @Description Returns a paginated list of loanUserCallRecordss starting after a given last id, useful for cursor-based pagination.
// @Tags loanUserCallRecords
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanUserCallRecordssReply{}
// @Router /api/v1/loanUserCallRecords/list [get]
// @Security BearerAuth
func (h *loanUserCallRecordsHandler) ListByLastID(c *gin.Context) {
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
	loanUserCallRecordss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanUserCallRecordss(loanUserCallRecordss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanUserCallRecords)
		return
	}

	response.Success(c, gin.H{
		"loanUserCallRecordss": data,
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
