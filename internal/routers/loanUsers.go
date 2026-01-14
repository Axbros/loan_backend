package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"

	"loan/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		loanUsersRouter(group, handler.NewLoanUsersHandler())
	})
}

func loanUsersRouter(group *gin.RouterGroup, h handler.LoanUsersHandler) {
	g := group.Group("/user")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	g.GET("/me", middleware.Auth(), h.Me)
	g.GET("/refer", h.Refer)
	g.POST("/", h.Create)          // [post] /api/v1/loanUsers
	g.DELETE("/:id", h.DeleteByID) // [delete] /api/v1/loanUsers/:id
	g.PUT("/:id", h.UpdateByID)    // [put] /api/v1/loanUsers/:id
	g.GET("/:id", h.GetByID)       // [get] /api/v1/loanUsers/:id
	g.POST("/list", h.List)        // [post] /api/v1/loanUsers/list

	g.POST("/delete/ids", h.DeleteByIDs)   // [post] /api/v1/loanUsers/delete/ids
	g.POST("/condition", h.GetByCondition) // [post] /api/v1/loanUsers/condition
	g.POST("/list/ids", h.ListByIDs)       // [post] /api/v1/loanUsers/list/ids
	g.GET("/list", h.ListByLastID)         // [get] /api/v1/loanUsers/list
}
