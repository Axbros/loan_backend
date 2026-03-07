package authz

import (
	"strconv"

	"github.com/gin-gonic/gin"
	smiddleware "github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/jwt"

	"loan/internal/cache"
	"loan/internal/dao"
	"loan/internal/database"
	"loan/internal/ecode"
)

func RequirePerm(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("claims")
		if !ok || v == nil {
			response.Out(c, ecode.Unauthorized)
			c.Abort()
			return
		}
		claims, ok := v.(*jwt.Claims)
		if !ok || claims == nil || claims.UID == "" {
			response.Out(c, ecode.Unauthorized)
			c.Abort()
			return
		}
		uid, err := strconv.ParseUint(claims.UID, 10, 64)
		if err != nil || uid == 0 {
			response.Out(c, ecode.Unauthorized)
			c.Abort()
			return
		}

		ctx := smiddleware.WrapCtx(c)
		usersDao := dao.NewLoanUsersDao(database.GetDB(), cache.NewLoanUsersCache(database.GetCacheType()))
		perms, err := usersDao.GetPermissionCodesByUserID(ctx, uid)
		if err != nil {
			response.Out(c, ecode.InternalServerError)
			c.Abort()
			return
		}
		if len(perms) == 0 {
			response.Out(c, ecode.Forbidden)
			c.Abort()
			return
		}
		has := false
		for _, p := range perms {
			if p == code {
				has = true
				break
			}
		}
		if !has {
			response.Out(c, ecode.Forbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
