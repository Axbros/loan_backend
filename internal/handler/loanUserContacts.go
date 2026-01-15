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

var _ LoanUserContactsHandler = (*loanUserContactsHandler)(nil)

// LoanUserContactsHandler defining the handler interface
type LoanUserContactsHandler interface {
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

type loanUserContactsHandler struct {
	iDao dao.LoanUserContactsDao
}

// NewLoanUserContactsHandler creating the handler interface
func NewLoanUserContactsHandler() LoanUserContactsHandler {
	return &loanUserContactsHandler{
		iDao: dao.NewLoanUserContactsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanUserContactsCache(database.GetCacheType()),
		),
	}
}

// Create a new loanUserContacts
// @Summary Create a new loanUserContacts
// @Description Creates a new loanUserContacts entity using the provided data in the request body.
// @Tags loanUserContacts
// @Accept json
// @Produce json
// @Param data body types.CreateLoanUserContactsRequest true "loanUserContacts information"
// @Success 200 {object} types.CreateLoanUserContactsReply{}
// @Router /api/v1/loanUserContacts [post]
// @Security BearerAuth
func (h *loanUserContactsHandler) Create(c *gin.Context) {
	form := &types.CreateLoanUserContactsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanUserContacts := &model.LoanUserContacts{}
	err = copier.Copy(loanUserContacts, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanUserContacts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanUserContacts)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanUserContacts.ID})
}

// DeleteByID delete a loanUserContacts by id
// @Summary Delete a loanUserContacts by id
// @Description Deletes a existing loanUserContacts identified by the given id in the path.
// @Tags loanUserContacts
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanUserContactsByIDReply{}
// @Router /api/v1/loanUserContacts/{id} [delete]
// @Security BearerAuth
func (h *loanUserContactsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanUserContactsIDFromPath(c)
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

// UpdateByID update a loanUserContacts by id
// @Summary Update a loanUserContacts by id
// @Description Updates the specified loanUserContacts by given id in the path, support partial update.
// @Tags loanUserContacts
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanUserContactsByIDRequest true "loanUserContacts information"
// @Success 200 {object} types.UpdateLoanUserContactsByIDReply{}
// @Router /api/v1/loanUserContacts/{id} [put]
// @Security BearerAuth
func (h *loanUserContactsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanUserContactsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanUserContactsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanUserContacts := &model.LoanUserContacts{}
	err = copier.Copy(loanUserContacts, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanUserContacts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanUserContacts)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanUserContacts by id
// @Summary Get a loanUserContacts by id
// @Description Gets detailed information of a loanUserContacts specified by the given id in the path.
// @Tags loanUserContacts
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanUserContactsByIDReply{}
// @Router /api/v1/loanUserContacts/{id} [get]
// @Security BearerAuth
func (h *loanUserContactsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanUserContactsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUserContacts, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanUserContactsObjDetail{}
	err = copier.Copy(data, loanUserContacts)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanUserContacts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanUserContacts": data})
}

// List get a paginated list of loanUserContactss by custom conditions
// @Summary Get a paginated list of loanUserContactss by custom conditions
// @Description Returns a paginated list of loanUserContacts based on query filters, including page number and size.
// @Tags loanUserContacts
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanUserContactssReply{}
// @Router /api/v1/loanUserContacts/list [post]
// @Security BearerAuth
func (h *loanUserContactsHandler) List(c *gin.Context) {
	form := &types.ListLoanUserContactssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUserContactss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanUserContactss(loanUserContactss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanUserContacts)
		return
	}

	response.Success(c, gin.H{
		"records": data,
		"total":   total,
	})
}

// DeleteByIDs batch delete loanUserContacts by ids
// @Summary Batch delete loanUserContacts by ids
// @Description Deletes multiple loanUserContacts by a list of id
// @Tags loanUserContacts
// @Param data body types.DeleteLoanUserContactssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanUserContactssByIDsReply{}
// @Router /api/v1/loanUserContacts/delete/ids [post]
// @Security BearerAuth
func (h *loanUserContactsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanUserContactssByIDsRequest{}
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

// GetByCondition get a loanUserContacts by custom condition
// @Summary Get a loanUserContacts by custom condition
// @Description Returns a single loanUserContacts that matches the specified filter conditions.
// @Tags loanUserContacts
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanUserContactsByConditionReply{}
// @Router /api/v1/loanUserContacts/condition [post]
// @Security BearerAuth
func (h *loanUserContactsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanUserContactsByConditionRequest{}
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
	loanUserContacts, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanUserContactsObjDetail{}
	err = copier.Copy(data, loanUserContacts)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanUserContacts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanUserContacts": data})
}

// ListByIDs batch get loanUserContacts by ids
// @Summary Batch get loanUserContacts by ids
// @Description Returns a list of loanUserContacts that match the list of id.
// @Tags loanUserContacts
// @Param data body types.ListLoanUserContactssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanUserContactssByIDsReply{}
// @Router /api/v1/loanUserContacts/list/ids [post]
// @Security BearerAuth
func (h *loanUserContactsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanUserContactssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanUserContactsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanUserContactss := []*types.LoanUserContactsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanUserContactsMap[id]; ok {
			record, err := convertLoanUserContacts(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanUserContacts)
				return
			}
			loanUserContactss = append(loanUserContactss, record)
		}
	}

	response.Success(c, gin.H{
		"loanUserContactss": loanUserContactss,
	})
}

// ListByLastID get a paginated list of loanUserContactss by last id
// @Summary Get a paginated list of loanUserContactss by last id
// @Description Returns a paginated list of loanUserContactss starting after a given last id, useful for cursor-based pagination.
// @Tags loanUserContacts
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanUserContactssReply{}
// @Router /api/v1/loanUserContacts/list [get]
// @Security BearerAuth
func (h *loanUserContactsHandler) ListByLastID(c *gin.Context) {
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
	loanUserContactss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanUserContactss(loanUserContactss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanUserContacts)
		return
	}

	response.Success(c, gin.H{
		"loanUserContactss": data,
	})
}

func getLoanUserContactsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanUserContacts(loanUserContacts *model.LoanUserContacts) (*types.LoanUserContactsObjDetail, error) {
	data := &types.LoanUserContactsObjDetail{}
	err := copier.Copy(data, loanUserContacts)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanUserContactss(fromValues []*model.LoanUserContacts) ([]*types.LoanUserContactsObjDetail, error) {
	toValues := []*types.LoanUserContactsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanUserContacts(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
