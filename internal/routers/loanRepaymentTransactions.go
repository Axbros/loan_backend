package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"

	"loan/internal/authz"
	"loan/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		loanRepaymentTransactionsRouter(group, handler.NewLoanRepaymentTransactionsHandler())
	})
}

//TODO 用户添加回款 如果大于等于应还金额应该把loanCollectionCases.go的status设置为已完成2

func loanRepaymentTransactionsRouter(group *gin.RouterGroup, h handler.LoanRepaymentTransactionsHandler) {
	g := group.Group("/repayment-transaction")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", authz.RequirePerm("repayment-transaction:add"), h.Create)             // [post] /api/v1/loanRepaymentTransactions
	g.DELETE("/:id", authz.RequirePerm("repayment-transaction:delete"), h.DeleteByID) // [delete] /api/v1/loanRepaymentTransactions/:id
	g.PUT("/:id", authz.RequirePerm("repayment-transaction:update"), h.UpdateByID)    // [put] /api/v1/loanRepaymentTransactions/:id
	g.GET("/:id", authz.RequirePerm("repayment-transaction:view"), h.GetByID)         // [get] /api/v1/loanRepaymentTransactions/:id
	g.POST("/list", authz.RequirePerm("repayment-transaction:view"), h.List)          // [post] /api/v1/loanRepaymentTransactions/list

	g.POST("/loan-info", authz.RequirePerm("repayment-transaction:view"), h.DetailByScheduleID)
	g.POST("/history", authz.RequirePerm("repayment-transaction:view"), h.History)
	g.POST("/upload-voucher", authz.RequirePerm("repayment-transaction:upload"), h.UploadVoucher)
	g.GET("/upload-voucher/:file_name", authz.RequirePerm("repayment-transaction:view"), h.GetVoucherBase64)
}
