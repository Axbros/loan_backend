package routers

import (
	"github.com/gin-gonic/gin"

	"loan/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		loanUserDeviceAppsRouter(group, handler.NewLoanUserDeviceAppsHandler())
	})
}

func loanUserDeviceAppsRouter(group *gin.RouterGroup, h handler.LoanUserDeviceAppsHandler) {
	g := group.Group("/loanUserDeviceApps")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", h.Create)          // [post] /api/v1/loanUserDeviceApps
	g.DELETE("/:id", h.DeleteByID) // [delete] /api/v1/loanUserDeviceApps/:id
	g.PUT("/:id", h.UpdateByID)    // [put] /api/v1/loanUserDeviceApps/:id
	g.GET("/:id", h.GetByID)       // [get] /api/v1/loanUserDeviceApps/:id
	g.POST("/list", h.List)        // [post] /api/v1/loanUserDeviceApps/list

	g.POST("/delete/ids", h.DeleteByIDs)   // [post] /api/v1/loanUserDeviceApps/delete/ids
	g.POST("/condition", h.GetByCondition) // [post] /api/v1/loanUserDeviceApps/condition
	g.POST("/list/ids", h.ListByIDs)       // [post] /api/v1/loanUserDeviceApps/list/ids
	g.GET("/list", h.ListByLastID)         // [get] /api/v1/loanUserDeviceApps/list
}
