package routers

import (
	"loan/internal/authz"

	"github.com/gin-gonic/gin"

	"loan/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		loanDepartmentsRouter(group, handler.NewLoanDepartmentsHandler())
	})
}

func loanDepartmentsRouter(group *gin.RouterGroup, h handler.LoanDepartmentsHandler) {
	g := group.Group("/department")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", authz.RequirePerm("department:add"), h.Create)             // [post] /api/v1/loanDepartments
	g.DELETE("/:id", authz.RequirePerm("department:delete"), h.DeleteByID) // [delete] /api/v1/loanDepartments/:id
	g.PUT("/:id", authz.RequirePerm("department:update"), h.UpdateByID)    // [put] /api/v1/loanDepartments/:id
	g.GET("/:id", h.GetByID)                                               // [get] /api/v1/loanDepartments/:id
	g.POST("/list", authz.RequirePerm("department:view"), h.List)          // [post] /api/v1/loanDepartments/list

}
