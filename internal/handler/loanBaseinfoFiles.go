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

var _ LoanBaseinfoFilesHandler = (*loanBaseinfoFilesHandler)(nil)

// LoanBaseinfoFilesHandler defining the handler interface
type LoanBaseinfoFilesHandler interface {
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

type loanBaseinfoFilesHandler struct {
	iDao dao.LoanBaseinfoFilesDao
}

// NewLoanBaseinfoFilesHandler creating the handler interface
func NewLoanBaseinfoFilesHandler() LoanBaseinfoFilesHandler {
	return &loanBaseinfoFilesHandler{
		iDao: dao.NewLoanBaseinfoFilesDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanBaseinfoFilesCache(database.GetCacheType()),
		),
	}
}

// Create a new loanBaseinfoFiles
// @Summary Create a new loanBaseinfoFiles
// @Description Creates a new loanBaseinfoFiles entity using the provided data in the request body.
// @Tags loanBaseinfoFiles
// @Accept json
// @Produce json
// @Param data body types.CreateLoanBaseinfoFilesRequest true "loanBaseinfoFiles information"
// @Success 200 {object} types.CreateLoanBaseinfoFilesReply{}
// @Router /api/v1/loanBaseinfoFiles [post]
// @Security BearerAuth
func (h *loanBaseinfoFilesHandler) Create(c *gin.Context) {
	form := &types.CreateLoanBaseinfoFilesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanBaseinfoFiles := &model.LoanBaseinfoFiles{}
	err = copier.Copy(loanBaseinfoFiles, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanBaseinfoFiles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanBaseinfoFiles)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanBaseinfoFiles.ID})
}

// DeleteByID delete a loanBaseinfoFiles by id
// @Summary Delete a loanBaseinfoFiles by id
// @Description Deletes a existing loanBaseinfoFiles identified by the given id in the path.
// @Tags loanBaseinfoFiles
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanBaseinfoFilesByIDReply{}
// @Router /api/v1/loanBaseinfoFiles/{id} [delete]
// @Security BearerAuth
func (h *loanBaseinfoFilesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanBaseinfoFilesIDFromPath(c)
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

// UpdateByID update a loanBaseinfoFiles by id
// @Summary Update a loanBaseinfoFiles by id
// @Description Updates the specified loanBaseinfoFiles by given id in the path, support partial update.
// @Tags loanBaseinfoFiles
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanBaseinfoFilesByIDRequest true "loanBaseinfoFiles information"
// @Success 200 {object} types.UpdateLoanBaseinfoFilesByIDReply{}
// @Router /api/v1/loanBaseinfoFiles/{id} [put]
// @Security BearerAuth
func (h *loanBaseinfoFilesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanBaseinfoFilesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanBaseinfoFilesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanBaseinfoFiles := &model.LoanBaseinfoFiles{}
	err = copier.Copy(loanBaseinfoFiles, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanBaseinfoFiles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanBaseinfoFiles)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanBaseinfoFiles by id
// @Summary Get a loanBaseinfoFiles by id
// @Description Gets detailed information of a loanBaseinfoFiles specified by the given id in the path.
// @Tags loanBaseinfoFiles
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanBaseinfoFilesByIDReply{}
// @Router /api/v1/loanBaseinfoFiles/{id} [get]
// @Security BearerAuth
func (h *loanBaseinfoFilesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanBaseinfoFilesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanBaseinfoFiles, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanBaseinfoFilesObjDetail{}
	err = copier.Copy(data, loanBaseinfoFiles)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanBaseinfoFiles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanBaseinfoFiles": data})
}

// List get a paginated list of loanBaseinfoFiless by custom conditions
// @Summary Get a paginated list of loanBaseinfoFiless by custom conditions
// @Description Returns a paginated list of loanBaseinfoFiles based on query filters, including page number and size.
// @Tags loanBaseinfoFiles
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanBaseinfoFilessReply{}
// @Router /api/v1/loanBaseinfoFiles/list [post]
// @Security BearerAuth
func (h *loanBaseinfoFilesHandler) List(c *gin.Context) {
	form := &types.ListLoanBaseinfoFilessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanBaseinfoFiless, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanBaseinfoFiless(loanBaseinfoFiless)
	if err != nil {
		response.Error(c, ecode.ErrListLoanBaseinfoFiles)
		return
	}

	response.Success(c, gin.H{
		"loanBaseinfoFiless": data,
		"total":              total,
	})
}

// DeleteByIDs batch delete loanBaseinfoFiles by ids
// @Summary Batch delete loanBaseinfoFiles by ids
// @Description Deletes multiple loanBaseinfoFiles by a list of id
// @Tags loanBaseinfoFiles
// @Param data body types.DeleteLoanBaseinfoFilessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanBaseinfoFilessByIDsReply{}
// @Router /api/v1/loanBaseinfoFiles/delete/ids [post]
// @Security BearerAuth
func (h *loanBaseinfoFilesHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanBaseinfoFilessByIDsRequest{}
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

// GetByCondition get a loanBaseinfoFiles by custom condition
// @Summary Get a loanBaseinfoFiles by custom condition
// @Description Returns a single loanBaseinfoFiles that matches the specified filter conditions.
// @Tags loanBaseinfoFiles
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanBaseinfoFilesByConditionReply{}
// @Router /api/v1/loanBaseinfoFiles/condition [post]
// @Security BearerAuth
func (h *loanBaseinfoFilesHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanBaseinfoFilesByConditionRequest{}
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
	loanBaseinfoFiles, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanBaseinfoFilesObjDetail{}
	err = copier.Copy(data, loanBaseinfoFiles)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanBaseinfoFiles)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanBaseinfoFiles": data})
}

// ListByIDs batch get loanBaseinfoFiles by ids
// @Summary Batch get loanBaseinfoFiles by ids
// @Description Returns a list of loanBaseinfoFiles that match the list of id.
// @Tags loanBaseinfoFiles
// @Param data body types.ListLoanBaseinfoFilessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanBaseinfoFilessByIDsReply{}
// @Router /api/v1/loanBaseinfoFiles/list/ids [post]
// @Security BearerAuth
func (h *loanBaseinfoFilesHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanBaseinfoFilessByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanBaseinfoFilesMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanBaseinfoFiless := []*types.LoanBaseinfoFilesObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanBaseinfoFilesMap[id]; ok {
			record, err := convertLoanBaseinfoFiles(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanBaseinfoFiles)
				return
			}
			loanBaseinfoFiless = append(loanBaseinfoFiless, record)
		}
	}

	response.Success(c, gin.H{
		"loanBaseinfoFiless": loanBaseinfoFiless,
	})
}

// ListByLastID get a paginated list of loanBaseinfoFiless by last id
// @Summary Get a paginated list of loanBaseinfoFiless by last id
// @Description Returns a paginated list of loanBaseinfoFiless starting after a given last id, useful for cursor-based pagination.
// @Tags loanBaseinfoFiles
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanBaseinfoFilessReply{}
// @Router /api/v1/loanBaseinfoFiles/list [get]
// @Security BearerAuth
func (h *loanBaseinfoFilesHandler) ListByLastID(c *gin.Context) {
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
	loanBaseinfoFiless, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanBaseinfoFiless(loanBaseinfoFiless)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanBaseinfoFiles)
		return
	}

	response.Success(c, gin.H{
		"loanBaseinfoFiless": data,
	})
}

func getLoanBaseinfoFilesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanBaseinfoFiles(loanBaseinfoFiles *model.LoanBaseinfoFiles) (*types.LoanBaseinfoFilesObjDetail, error) {
	data := &types.LoanBaseinfoFilesObjDetail{}
	err := copier.Copy(data, loanBaseinfoFiles)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanBaseinfoFiless(fromValues []*model.LoanBaseinfoFiles) ([]*types.LoanBaseinfoFilesObjDetail, error) {
	toValues := []*types.LoanBaseinfoFilesObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanBaseinfoFiles(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
