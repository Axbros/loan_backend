package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"

	"loan/internal/authz"
	"loan/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		loanRiskCustomerRouter(group, handler.NewLoanRiskCustomerHandler())
	})
}

func loanRiskCustomerRouter(group *gin.RouterGroup, h handler.LoanRiskCustomerHandler) {
	g := group.Group("/risk_customer")

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", authz.RequirePerm("risk-customer:add"), h.Create)             // [post] /api/v1/loanRiskCustomer
	g.DELETE("/:id", authz.RequirePerm("risk-customer:delete"), h.DeleteByID) // [delete] /api/v1/loanRiskCustomer/:id
	g.PUT("/:id", authz.RequirePerm("risk-customer:update"), h.UpdateByID)    // [put] /api/v1/loanRiskCustomer/:id
	g.GET("/:id", authz.RequirePerm("risk-customer:view"), h.GetByID)         // [get] /api/v1/loanRiskCustomer/:id
	g.POST("/list", authz.RequirePerm("risk-customer:view"), h.List)          // [post] /api/v1/loanRiskCustomer/list
}
