package routers

import (
	"github.com/gin-gonic/gin"

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
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", h.Create)          // [post] /api/v1/loanBaseinfo
	g.DELETE("/:id", h.DeleteByID) // [delete] /api/v1/loanBaseinfo/:id
	g.PUT("/:id", h.UpdateByID)    // [put] /api/v1/loanBaseinfo/:id
	g.GET("/:id", h.GetByID)       // [get] /api/v1/loanBaseinfo/:id
	g.POST("/list", h.List)        // [post] /api/v1/loanBaseinfo/list

	g.POST("/delete/ids", h.DeleteByIDs)   // [post] /api/v1/loanBaseinfo/delete/ids
	g.POST("/condition", h.GetByCondition) // [post] /api/v1/loanBaseinfo/condition
	g.POST("/list/ids", h.ListByIDs)       // [post] /api/v1/loanBaseinfo/list/ids
	g.GET("/list", h.ListByLastID)         // [get] /api/v1/loanBaseinfo/list
}
