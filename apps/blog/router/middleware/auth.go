package middleware

import (
	"fmt"
	"github.com/shipengqi/example.v1/apps/blog/pkg/app"
	"github.com/shipengqi/example.v1/apps/blog/pkg/e"
	"github.com/shipengqi/example.v1/apps/blog/service"

	"github.com/gin-gonic/gin"
)

func Authenticate(s *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip for the login request.
		// path := c.Request.URL.Path
		// if path == "/api/v1/users" {
		// 	   c.Next()
		// 	   return
		// }

		var token string
		authorization := c.GetHeader("Authorization")
		xToken := c.GetHeader("X-AUTH-TOKEN")
		if len(xToken) > 0 {
			token = xToken
		} else if len(authorization) > 0 {
			// get the token part
			_, _ = fmt.Sscanf(authorization, "Bearer %s", &token)
		} else {
			if t, ok := c.GetQuery("token"); ok {
				token = t
			}
		}

		if len(token) == 0 {
			app.SendResponse(c, e.ErrUnauthorized, nil)
			c.Abort()
			return
		}

		claims, err := s.AuthSvc.Authenticate(token)
		if err != nil {
			app.SendResponse(c, err, nil)
			c.Abort()
			return
		}
		c.Set("auth_claims", claims)

		err = s.AuthSvc.Authorize(claims, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			app.SendResponse(c, err, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
