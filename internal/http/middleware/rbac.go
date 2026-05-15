package middleware

import (
	"net/http"

	"driving-authority-backend/internal/domain"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...domain.Role) gin.HandlerFunc {
	allowed := map[domain.Role]struct{}{}
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *gin.Context) {
		u, ok := c.Get(authUserKey)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		user := u.(AuthUser)
		if _, ok := allowed[user.Role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
