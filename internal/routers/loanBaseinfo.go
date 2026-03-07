package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"

	"loan/internal/authz"
	"loan/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		loanBaseinfoRouter(group, handler.NewLoanBaseinfoHandler())
	})
}

func loanBaseinfoRouter(group *gin.RouterGroup, h handler.LoanBaseinfoHandler) {
	g := group.Group("/customer")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", authz.RequirePerm("customer:add"), h.Create)
	g.DELETE("/:id", authz.RequirePerm("customer:delete"), h.DeleteByID)
	g.PUT("/:id", authz.RequirePerm("customer:update"), h.UpdateByID)
	g.GET("/:id", authz.RequirePerm("customer:view"), h.GetByID)
	g.POST("/list", authz.RequirePerm("customer:view"), h.List)

	g.POST("/pre-review", authz.RequirePerm("loan:pre_review"), h.PreReview)
	g.POST("/finance-review", authz.RequirePerm("loan:finance_review"), h.FinanceReview)

	g.POST("/withAuditRecord/list", authz.RequirePerm("customer:view"), h.WithAuditRecordList)
}
