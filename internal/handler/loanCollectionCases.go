package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"
	"github.com/go-sql-driver/mysql"

	"loan/internal/cache"
	"loan/internal/dao"
	"loan/internal/database"
	"loan/internal/ecode"
	"loan/internal/model"
	"loan/internal/types"
)

var _ LoanCollectionCasesHandler = (*loanCollectionCasesHandler)(nil)

// LoanCollectionCasesHandler defining the handler interface
type LoanCollectionCasesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Assign(c *gin.Context)
}

type loanCollectionCasesHandler struct {
	iDao        dao.LoanCollectionCasesDao
	scheduleDao dao.LoanRepaymentSchedulesDao
}

// NewLoanCollectionCasesHandler creating the handler interface
func NewLoanCollectionCasesHandler() LoanCollectionCasesHandler {
	return &loanCollectionCasesHandler{
		iDao: dao.NewLoanCollectionCasesDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanCollectionCasesCache(database.GetCacheType()),
		),
		scheduleDao: dao.NewLoanRepaymentSchedulesDao(
			database.GetDB(),
			cache.NewLoanRepaymentSchedulesCache(database.GetCacheType()),
		),
	}
}
func (h *loanCollectionCasesHandler) Assign(c *gin.Context) {
	form := &types.CreateLoanCollectionCasesAssignRequest{}
	if err := c.ShouldBindJSON(form); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	uid, ok := getUIDFromClaims(c)
	if !ok || uid == 0 {
		response.Error(c, ecode.Unauthorized)
		return
	}

	if form.CollectorUserID == 0 || len(form.ScheduleIDs) == 0 {
		response.Error(c, ecode.InvalidParams)
		return
	}

	uniq := make([]uint64, 0, len(form.ScheduleIDs))
	seen := make(map[uint64]struct{}, len(form.ScheduleIDs))
	for _, sid := range form.ScheduleIDs {
		if sid == 0 {
			continue
		}
		if _, exists := seen[sid]; exists {
			continue
		}
		seen[sid] = struct{}{}
		uniq = append(uniq, sid)
	}
	if len(uniq) == 0 {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)

	db := database.GetDB()
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.Error("Begin tx error", logger.Err(tx.Error), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InternalServerError)
		return
	}

	// 兜底：函数任何提前 return 都会回滚（Commit 后 Rollback 不会生效）
	defer func() {
		_ = tx.Rollback().Error
	}()

	duplicateIDs := make([]uint64, 0)

	for _, sid := range uniq {
		record := &model.LoanCollectionCases{
			ScheduleID:       sid,
			CollectorUserID:  form.CollectorUserID,
			AssignedByUserID: uid,
		}

		_, err := h.iDao.CreateByTx(ctx, tx, record)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				// 重复指派，跳过
				duplicateIDs = append(duplicateIDs, sid)
				logger.Warn("重复指派，跳过 LoanCollectionCases CreateByTx error: ", logger.Err(err))
				continue
			}
			logger.Error(
				"CreateByTx error",
				logger.Err(err),
				logger.Any("record", record),
				middleware.GCtxRequestIDField(c),
			)

			_ = tx.Rollback().Error
			response.Error(c, ecode.InternalServerError)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Error("Commit tx error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InternalServerError)
		return
	}

	response.Success(c, gin.H{
		"duplicate_schedule_ids":   duplicateIDs,
		"duplicate_schedule_total": len(duplicateIDs),
		"success_assign_totla":     len(form.ScheduleIDs) - len(duplicateIDs),
	})
}

// Create a new loanCollectionCases
// @Summary Create a new loanCollectionCases
// @Description Creates a new loanCollectionCases entity using the provided data in the request body.
// @Tags loanCollectionCases
// @Accept json
// @Produce json
// @Param data body types.CreateLoanCollectionCasesRequest true "loanCollectionCases information"
// @Success 200 {object} types.CreateLoanCollectionCasesReply{}
// @Router /api/v1/loanCollectionCases [post]
// @Security BearerAuth
func (h *loanCollectionCasesHandler) Create(c *gin.Context) {
	form := &types.CreateLoanCollectionCasesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	loanCollectionCases := &model.LoanCollectionCases{}
	err = copier.Copy(loanCollectionCases, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLoanCollectionCases)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, loanCollectionCases)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": loanCollectionCases.ID})
}

// DeleteByID delete a loanCollectionCases by id
// @Summary Delete a loanCollectionCases by id
// @Description Deletes a existing loanCollectionCases identified by the given id in the path.
// @Tags loanCollectionCases
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLoanCollectionCasesByIDReply{}
// @Router /api/v1/loanCollectionCases/{id} [delete]
// @Security BearerAuth
func (h *loanCollectionCasesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLoanCollectionCasesIDFromPath(c)
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

// UpdateByID update a loanCollectionCases by id
// @Summary Update a loanCollectionCases by id
// @Description Updates the specified loanCollectionCases by given id in the path, support partial update.
// @Tags loanCollectionCases
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLoanCollectionCasesByIDRequest true "loanCollectionCases information"
// @Success 200 {object} types.UpdateLoanCollectionCasesByIDReply{}
// @Router /api/v1/loanCollectionCases/{id} [put]
// @Security BearerAuth
func (h *loanCollectionCasesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLoanCollectionCasesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLoanCollectionCasesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	loanCollectionCases := &model.LoanCollectionCases{}
	err = copier.Copy(loanCollectionCases, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLoanCollectionCases)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, loanCollectionCases)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a loanCollectionCases by id
// @Summary Get a loanCollectionCases by id
// @Description Gets detailed information of a loanCollectionCases specified by the given id in the path.
// @Tags loanCollectionCases
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLoanCollectionCasesByIDReply{}
// @Router /api/v1/loanCollectionCases/{id} [get]
// @Security BearerAuth
func (h *loanCollectionCasesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLoanCollectionCasesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanCollectionCases, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LoanCollectionCasesObjDetail{}
	err = copier.Copy(data, loanCollectionCases)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLoanCollectionCases)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"loanCollectionCases": data})
}

// List get a paginated list of loanCollectionCasess by custom conditions
// @Summary Get a paginated list of loanCollectionCasess by custom conditions
// @Description Returns a paginated list of loanCollectionCases based on query filters, including page number and size.
// @Tags loanCollectionCases
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLoanCollectionCasessReply{}
// @Router /api/v1/loanCollectionCases/list [post]
// @Security BearerAuth
func (h *loanCollectionCasesHandler) List(c *gin.Context) {
	form := &types.ListLoanCollectionCasessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	loanCollectionCasess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{
		"records": loanCollectionCasess,
		"total":   total,
	})
}

func getLoanCollectionCasesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLoanCollectionCases(loanCollectionCases *model.LoanCollectionCases) (*types.LoanCollectionCasesObjDetail, error) {
	data := &types.LoanCollectionCasesObjDetail{}
	err := copier.Copy(data, loanCollectionCases)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLoanCollectionCasess(fromValues []*model.LoanCollectionCases) ([]*types.LoanCollectionCasesObjDetail, error) {
	toValues := []*types.LoanCollectionCasesObjDetail{}
	for _, v := range fromValues {
		data, err := convertLoanCollectionCases(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
