package handler

import (
	"github.com/gin-gonic/gin"

	"loan/internal/cache"
	"loan/internal/dao"
	"loan/internal/database"
	"loan/internal/ecode"
	"loan/internal/types"

	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
)

var _ LoanAuditsHandler = (*loanAuditsHandler)(nil)

// LoanAuditsHandler defining the handler interface
type LoanAuditsHandler interface {
	Detail(c *gin.Context)
}

type loanAuditsHandler struct {
	iDao dao.LoanAuditsDao
}

// NewLoanAuditsHandler creating the handler interface
func NewLoanAuditsHandler() LoanAuditsHandler {
	return &loanAuditsHandler{
		iDao: dao.NewLoanAuditsDao(
			database.GetDB(), // db driver is mysql
			cache.NewLoanAuditsCache(database.GetCacheType()),
		),
	}
}

func (h *loanAuditsHandler) Detail(c *gin.Context) {
	form := &types.RequestLoanAuditsDetail{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)
	record, err := h.iDao.GetByBaseinfoID(ctx, form.BaseinfoID, form.AuditType)
	if err != nil {
		logger.Warn("GetByBaseinfoID error: ", logger.Err(err))
		response.Error(c, ecode.InternalServerError)
		return
	}
	var disbursmentRecord *types.DisbursementWithChannel // 核心修正：改为放款记录模型类型

	if form.AuditType == FinanceReviewType && record != nil && record.AuditResult == 1 {

		disbursmentRecord, err = h.iDao.GetDisbursmentsByBaseInfoID(c, form.BaseinfoID)
		if err != nil {
			logger.Warn("GetByCondition error: ", logger.Err(err))
			response.Error(c, ecode.ErrGetByConditionLoanDisbursements)
			return
		}
	}
	response.Success(c, gin.H{
		"record": record,
		"extra":  disbursmentRecord,
	})

}
