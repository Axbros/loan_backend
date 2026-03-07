package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"

	"loan/internal/authz"
	"loan/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		loanReferralVisitsRouter(group, handler.NewLoanReferralVisitsHandler())
	})
}

func loanReferralVisitsRouter(group *gin.RouterGroup, h handler.LoanReferralVisitsHandler) {
	g := group.Group("/loanReferralVisits")

	// JWT authentication reference: https://go-sponge.com/component/transport/gin.html#jwt-authorization-middleware

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", authz.RequirePerm("referral-visit:add"), h.Create)             // [post] /api/v1/loanReferralVisits
	g.DELETE("/:id", authz.RequirePerm("referral-visit:delete"), h.DeleteByID) // [delete] /api/v1/loanReferralVisits/:id
	g.PUT("/:id", authz.RequirePerm("referral-visit:update"), h.UpdateByID)    // [put] /api/v1/loanReferralVisits/:id
	g.GET("/:id", authz.RequirePerm("referral-visit:view"), h.GetByID)         // [get] /api/v1/loanReferralVisits/:id
	g.POST("/list", authz.RequirePerm("referral-visit:view"), h.List)          // [post] /api/v1/loanReferralVisits/list

}
