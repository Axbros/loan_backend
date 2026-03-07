package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"

	"loan/internal/authz"
	"loan/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		loanRolesRouter(group, handler.NewLoanRolesHandler())
	})
}

func loanRolesRouter(group *gin.RouterGroup, h handler.LoanRolesHandler) {
	g := group.Group("/roles")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", authz.RequirePerm("role:add"), h.Create)             // [post] /api/v1/loanRoles
	g.DELETE("/:id", authz.RequirePerm("role:delete"), h.DeleteByID) // [delete] /api/v1/loanRoles/:id
	g.PUT("/:id", authz.RequirePerm("role:update"), h.UpdateByID)    // [put] /api/v1/loanRoles/:id
	g.GET("/:id", authz.RequirePerm("role:view"), h.GetByID)         // [get] /api/v1/loanRoles/:id
	g.POST("/list", authz.RequirePerm("role:view"), h.List)          // [post] /api/v1/loanRoles/list

}
