package routers

import (
	"loan/internal/authz"
	"loan/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
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

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", h.Create)
	g.DELETE("/:id", middleware.Auth(), authz.RequirePerm("customer:delete"), h.DeleteByID)
	g.PUT("/:id", middleware.Auth(), authz.RequirePerm("customer:update"), h.UpdateByID)
	g.GET("/:id", middleware.Auth(), authz.RequirePerm("customer:view"), h.GetByID)
	g.POST("/list", middleware.Auth(), authz.RequirePerm("customer:view"), h.List)

	g.POST("/pre-review", middleware.Auth(), authz.RequirePerm("loan:pre_review"), h.PreReview)
	g.POST("/finance-review", middleware.Auth(), authz.RequirePerm("loan:finance_review"), h.FinanceReview)

	g.POST("/withAuditRecord/list", middleware.Auth(), authz.RequirePerm("customer:view"), h.WithAuditRecordList)

	g.POST("/upload-certificate", h.UploadCertificate)
	g.GET("/upload-certificate/:file_name", authz.RequirePerm("customer:view"), h.GetCertificateBase64)
}
