package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"

	"loan/internal/authz"
	"loan/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		loanRepaymentSchedulesRouter(group, handler.NewLoanRepaymentSchedulesHandler())
	})
}

func loanRepaymentSchedulesRouter(group *gin.RouterGroup, h handler.LoanRepaymentSchedulesHandler) {
	g := group.Group("/repayment-schedule")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.
	g.POST("/overview", authz.RequirePerm("repayment-schedule:view"), h.Overview)
	g.POST("/", authz.RequirePerm("repayment-schedule:add"), h.Create)             // [post] /api/v1/loanRepaymentSchedules
	g.DELETE("/:id", authz.RequirePerm("repayment-schedule:delete"), h.DeleteByID) // [delete] /api/v1/loanRepaymentSchedules/:id
	g.PUT("/:id", authz.RequirePerm("repayment-schedule:update"), h.UpdateByID)    // [put] /api/v1/loanRepaymentSchedules/:id
	g.GET("/:id", authz.RequirePerm("repayment-schedule:view"), h.GetByID)         // [get] /api/v1/loanRepaymentSchedules/:id
	g.POST("/list", authz.RequirePerm("repayment-schedule:view"), h.List)          // [post] /api/v1/loanRepaymentSchedules/list

}
