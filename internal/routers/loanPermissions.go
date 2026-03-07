package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"

	"loan/internal/authz"
	"loan/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		loanPermissionsRouter(group, handler.NewLoanPermissionsHandler())
	})
}

func loanPermissionsRouter(group *gin.RouterGroup, h handler.LoanPermissionsHandler) {
	g := group.Group("/loanPermissions")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", authz.RequirePerm("permission:add"), h.Create)             // [post] /api/v1/loanPermissions
	g.DELETE("/:id", authz.RequirePerm("permission:delete"), h.DeleteByID) // [delete] /api/v1/loanPermissions/:id
	g.PUT("/:id", authz.RequirePerm("permission:update"), h.UpdateByID)    // [put] /api/v1/loanPermissions/:id
	g.GET("/:id", authz.RequirePerm("permission:view"), h.GetByID)         // [get] /api/v1/loanPermissions/:id
	g.POST("/list", authz.RequirePerm("permission:view"), h.List)          // [post] /api/v1/loanPermissions/list
	g.POST("/permission_update", authz.RequirePerm("permission:update"), h.PermissionUpdate)

}
