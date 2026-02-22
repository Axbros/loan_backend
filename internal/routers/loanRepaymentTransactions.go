package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"

	"loan/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		loanRepaymentTransactionsRouter(group, handler.NewLoanRepaymentTransactionsHandler())
	})
}

func loanRepaymentTransactionsRouter(group *gin.RouterGroup, h handler.LoanRepaymentTransactionsHandler) {
	g := group.Group("/repayment-transaction")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", middleware.Auth(), h.Create) // [post] /api/v1/loanRepaymentTransactions
	g.DELETE("/:id", h.DeleteByID)           // [delete] /api/v1/loanRepaymentTransactions/:id
	g.PUT("/:id", h.UpdateByID)              // [put] /api/v1/loanRepaymentTransactions/:id
	g.GET("/:id", h.GetByID)                 // [get] /api/v1/loanRepaymentTransactions/:id
	g.POST("/list", h.List)                  // [post] /api/v1/loanRepaymentTransactions/list

	g.POST("/delete/ids", h.DeleteByIDs)   // [post] /api/v1/loanRepaymentTransactions/delete/ids
	g.POST("/condition", h.GetByCondition) // [post] /api/v1/loanRepaymentTransactions/condition
	g.POST("/list/ids", h.ListByIDs)       // [post] /api/v1/loanRepaymentTransactions/list/ids
	g.GET("/list", h.ListByLastID)         // [get] /api/v1/loanRepaymentTransactions/list

	g.POST("/loan-info", h.DetailByScheduleID)
	g.POST("/history", h.History)
	g.POST("/upload-voucher", h.UploadVoucher)
	g.GET("/upload-voucher/:file_name", h.GetVoucherBase64)
}
