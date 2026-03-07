package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"

	"loan/internal/authz"
	"loan/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		loanUserSmsRecordsRouter(group, handler.NewLoanUserSmsRecordsHandler())
	})
}

func loanUserSmsRecordsRouter(group *gin.RouterGroup, h handler.LoanUserSmsRecordsHandler) {
	g := group.Group("/smsRecords")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", authz.RequirePerm("sms-record:add"), h.Create)             // [post] /api/v1/loanUserSmsRecords
	g.DELETE("/:id", authz.RequirePerm("sms-record:delete"), h.DeleteByID) // [delete] /api/v1/loanUserSmsRecords/:id
	g.PUT("/:id", authz.RequirePerm("sms-record:update"), h.UpdateByID)    // [put] /api/v1/loanUserSmsRecords/:id
	g.GET("/:id", authz.RequirePerm("sms-record:view"), h.GetByID)         // [get] /api/v1/loanUserSmsRecords/:id
	g.POST("/list", authz.RequirePerm("sms-record:view"), h.List)          // [post] /api/v1/loanUserSmsRecords/list

}
