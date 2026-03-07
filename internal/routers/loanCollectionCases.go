package routers

import (
	"loan/internal/handler"

	"loan/internal/authz"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		loanCollectionCasesRouter(group, handler.NewLoanCollectionCasesHandler())
	})
}

func loanCollectionCasesRouter(group *gin.RouterGroup, h handler.LoanCollectionCasesHandler) {
	g := group.Group("/collection-cases")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", authz.RequirePerm("collection-cases:add"), h.Create)             // [post] /api/v1/loanCollectionCases
	g.DELETE("/:id", authz.RequirePerm("collection-cases:delete"), h.DeleteByID) // [delete] /api/v1/loanCollectionCases/:id
	g.PUT("/:id", authz.RequirePerm("collection-cases:update"), h.UpdateByID)    // [put] /api/v1/loanCollectionCases/:id
	g.GET("/:id", authz.RequirePerm("collection-cases:view"), h.GetByID)         // [get] /api/v1/loanCollectionCases/:id
	g.POST("/list", authz.RequirePerm("collection-cases:view"), h.List)          // [post] /api/v1/loanCollectionCases/list
	g.POST("/assign", authz.RequirePerm("collection-cases:assign"), h.Assign)

}
