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

var _ LoanDisbursementsHandler = (*loanDisbursementsHandler)(nil)

// LoanDisbursementsHandler defining the handler interface
type LoanDisbursementsHandler interface {
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

type loanDisbursementsHandler struct {
	iDao        dao.LoanDisbursementsDao
	baseinfoDao dao.LoanBaseinfoDao
	auditDao    dao.LoanAuditsDao
}

// NewLoanDisbursementsHandler creating the handler interface
func NewLoanDisbursementsHandler() LoanDisbursementsHandler {
	return &loanDisbursementsHandler{
		iDao: dao.NewLoanDisbursementsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanDisbursementsCache(database.GetCacheType()),
		),
		baseinfoDao: dao.NewLoanBaseinfoDao(
			database.GetDB(),
			cache.NewLoanBaseinfoCache(database.GetCacheType()),
		),
		auditDao: dao.NewLoanAuditsDao(database.GetDB(), cache.NewLoanAuditsCache(database.GetCacheType())),
	}
}

// Create a new loanDisbursements
// @Summary Create a new loanDisbursements
// @Description Creates a new loanDisbursements entity using the provided data in the request body.
// @Tags loanDisbursements
// @Accept json
// @Produce json
// @Param data body types.CreateLoanDisbursementsRequest true "loanDisbursements information"
// @Success 200 {object} types.CreateLoanDisbursementsReply{}
// @Router /api/v1/loanDisbursements [post]
// @Security BearerAuth
func (h *loanDisbursementsHandler) Create(c *gin.Context) {
	form := &types.CreateLoanDisbursementsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanDisbursements := &model.LoanDisbursements{}
	err = copier.Copy(loanDisbursements, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanDisbursements)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanDisbursements)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanDisbursements.ID})
}

// DeleteByID delete a loanDisbursements by id
// @Summary Delete a loanDisbursements by id
// @Description Deletes a existing loanDisbursements identified by the given id in the path.
// @Tags loanDisbursements
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanDisbursementsByIDReply{}
// @Router /api/v1/loanDisbursements/{id} [delete]
// @Security BearerAuth
func (h *loanDisbursementsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanDisbursementsIDFromPath(c)
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

// UpdateByID update a loanDisbursements by id
// @Summary Update a loanDisbursements by id
// @Description Updates the specified loanDisbursements by given id in the path, support partial update.
// @Tags loanDisbursements
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanDisbursementsByIDRequest true "loanDisbursements information"
// @Success 200 {object} types.UpdateLoanDisbursementsByIDReply{}
// @Router /api/v1/loanDisbursements/{id} [put]
// @Security BearerAuth
func (h *loanDisbursementsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanDisbursementsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanDisbursementsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanDisbursements := &model.LoanDisbursements{}
	err = copier.Copy(loanDisbursements, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanDisbursements)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanDisbursements)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanDisbursements by id
// @Summary Get a loanDisbursements by id
// @Description Gets detailed information of a loanDisbursements specified by the given id in the path.
// @Tags loanDisbursements
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanDisbursementsByIDReply{}
// @Router /api/v1/loanDisbursements/{id} [get]
// @Security BearerAuth
func (h *loanDisbursementsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanDisbursementsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)

	// 1) 放款单
	disb, err := h.iDao.GetByID(ctx, id)
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

	disbDTO := &types.LoanDisbursementsObjDetail{}
	if err := copier.Copy(disbDTO, disb); err != nil {
		response.Error(c, ecode.ErrGetByIDLoanDisbursements)
		return
	}

	// 2) baseinfo
	baseinfo, err := h.baseinfoDao.GetByID(ctx, uint64(disb.BaseinfoID))
	if err != nil {
		logger.Error("Get baseinfo error", logger.Err(err), logger.Any("baseinfo_id", disb.BaseinfoID), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	baseDTO := &types.LoanBaseinfoObjDetail{}
	if err := copier.Copy(baseDTO, baseinfo); err != nil {
		response.Error(c, ecode.ErrGetByIDLoanBaseinfo)
		return
	}

	// 3) files（按 type 聚合）
	files, err := h.baseinfoDao.GetFilesMapByBaseinfoID(ctx, uint64(disb.BaseinfoID))
	if err != nil {
		logger.Error("GetFilesMapByBaseinfoID error", logger.Err(err), logger.Any("baseinfo_id", disb.BaseinfoID), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	if files == nil {
		files = map[string][]string{}
	}
	baseDTO.Files = files

	// 4) audits
	auditRecords, err := h.auditDao.ListByBaseinfoID(ctx, uint64(disb.BaseinfoID))
	if err != nil {
		logger.Error("List audits error", logger.Err(err), logger.Any("baseinfo_id", disb.BaseinfoID), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	audits := make([]*types.LoanAuditsObjDetail, 0, len(auditRecords))
	for _, r := range auditRecords {
		a := &types.LoanAuditsObjDetail{}
		_ = copier.Copy(a, r)
		// 如果 copier 对某些字段不工作，在这里手动补
		audits = append(audits, a)
	}

	// 5) 返回聚合
	response.Success(c, gin.H{
		"disbursement": disbDTO,
		"baseinfo":     baseDTO,
		"audits":       audits,
	})
}

// List get a paginated list of loanDisbursementss by custom conditions
// @Summary Get a paginated list of loanDisbursementss by custom conditions
// @Description Returns a paginated list of loanDisbursements based on query filters, including page number and size.
// @Tags loanDisbursements
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanDisbursementssReply{}
// @Router /api/v1/loanDisbursements/list [post]
// @Security BearerAuth
func (h *loanDisbursementsHandler) List(c *gin.Context) {
	form := &types.ListLoanDisbursementssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanDisbursementss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanDisbursementss(loanDisbursementss)
	if err != nil {
		response.Error(c, ecode.ErrListLoanDisbursements)
		return
	}

	response.Success(c, gin.H{
		"loanDisbursementss": data,
		"total":              total,
	})
}

// DeleteByIDs batch delete loanDisbursements by ids
// @Summary Batch delete loanDisbursements by ids
// @Description Deletes multiple loanDisbursements by a list of id
// @Tags loanDisbursements
// @Param data body types.DeleteLoanDisbursementssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLoanDisbursementssByIDsReply{}
// @Router /api/v1/loanDisbursements/delete/ids [post]
// @Security BearerAuth
func (h *loanDisbursementsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLoanDisbursementssByIDsRequest{}
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

// GetByCondition get a loanDisbursements by custom condition
// @Summary Get a loanDisbursements by custom condition
// @Description Returns a single loanDisbursements that matches the specified filter conditions.
// @Tags loanDisbursements
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanDisbursementsByConditionReply{}
// @Router /api/v1/loanDisbursements/condition [post]
// @Security BearerAuth
func (h *loanDisbursementsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLoanDisbursementsByConditionRequest{}
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
	loanDisbursements, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LoanDisbursementsObjDetail{}
	err = copier.Copy(data, loanDisbursements)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanDisbursements)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanDisbursements": data})
}

// ListByIDs batch get loanDisbursements by ids
// @Summary Batch get loanDisbursements by ids
// @Description Returns a list of loanDisbursements that match the list of id.
// @Tags loanDisbursements
// @Param data body types.ListLoanDisbursementssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLoanDisbursementssByIDsReply{}
// @Router /api/v1/loanDisbursements/list/ids [post]
// @Security BearerAuth
func (h *loanDisbursementsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLoanDisbursementssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanDisbursementsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	loanDisbursementss := []*types.LoanDisbursementsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := loanDisbursementsMap[id]; ok {
			record, err := convertLoanDisbursements(v)
			if err != nil {
				response.Error(c, ecode.ErrListLoanDisbursements)
				return
			}
			loanDisbursementss = append(loanDisbursementss, record)
		}
	}

	response.Success(c, gin.H{
		"loanDisbursementss": loanDisbursementss,
	})
}

// ListByLastID get a paginated list of loanDisbursementss by last id
// @Summary Get a paginated list of loanDisbursementss by last id
// @Description Returns a paginated list of loanDisbursementss starting after a given last id, useful for cursor-based pagination.
// @Tags loanDisbursements
// @Accept json
// @Produce json
// @Param lastID query int false "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLoanDisbursementssReply{}
// @Router /api/v1/loanDisbursements/list [get]
// @Security BearerAuth
func (h *loanDisbursementsHandler) ListByLastID(c *gin.Context) {
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
	loanDisbursementss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("lastID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLoanDisbursementss(loanDisbursementss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLoanDisbursements)
		return
	}

	response.Success(c, gin.H{
		"loanDisbursementss": data,
	})
}

func getLoanDisbursementsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanDisbursements(loanDisbursements *model.LoanDisbursements) (*types.LoanDisbursementsObjDetail, error) {
	data := &types.LoanDisbursementsObjDetail{}
	err := copier.Copy(data, loanDisbursements)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanDisbursementss(fromValues []*model.LoanDisbursements) ([]*types.LoanDisbursementsObjDetail, error) {
	toValues := []*types.LoanDisbursementsObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanDisbursements(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
