package middleware

import (
	"clinic/config"
	"net/http"
	"strings"

	"clinic/response"

	"github.com/gin-gonic/gin"
)

func TokenAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized)
			c.Abort()
			return
		}
		parts := strings.SplitN(auth, "Bearer ", 2)
		if len(parts) != 2 || parts[1] != cfg.AuthToken {
			response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized)
			c.Abort()
			return
		}
		c.Next()
	}
}
