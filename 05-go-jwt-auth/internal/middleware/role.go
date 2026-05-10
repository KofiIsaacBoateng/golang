package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func OnlyAdmin() gin.HandlerFunc{
	return func(c *gin.Context) {
		role, ok := GetRole(c);
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ok": false,
				"error": "Unauthorized!",
			})
			return;
		}

		if !strings.EqualFold(role, "admin") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ok": false,
				"error": "Unauthorized user... Only admin allowed!",
			})
			return;
		}

		c.Next()
	}
}