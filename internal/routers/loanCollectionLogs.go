package routers

import (
	"loan/internal/handler"

	"loan/internal/authz"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		loanCollectionLogsRouter(group, handler.NewLoanCollectionLogsHandler())
	})
}

func loanCollectionLogsRouter(group *gin.RouterGroup, h handler.LoanCollectionLogsHandler) {
	g := group.Group("/collection-logs")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", authz.RequirePerm("collection-logs:add"), h.Create)             // [post] /api/v1/loanCollectionLogs
	g.DELETE("/:id", authz.RequirePerm("collection-logs:delete"), h.DeleteByID) // [delete] /api/v1/loanCollectionLogs/:id
	g.PUT("/:id", authz.RequirePerm("collection-logs:update"), h.UpdateByID)    // [put] /api/v1/loanCollectionLogs/:id
	g.GET("/:id", authz.RequirePerm("collection-logs:view"), h.GetByID)         // [get] /api/v1/loanCollectionLogs/:id
	g.POST("/list", authz.RequirePerm("collection-logs:view"), h.List)          // [post] /api/v1/loanCollectionLogs/list

}
