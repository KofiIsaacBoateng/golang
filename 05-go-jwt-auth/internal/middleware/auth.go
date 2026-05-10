package middleware

import (
	"jwt-auth/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ctxUserIdKey = "user.id"
	ctxUserRoleKey = "user.role" 
)

func Authorized(jwtSecret string) gin.HandlerFunc {
	return func (c *gin.Context) {
		// get auth header
		authheader := strings.TrimSpace(c.GetHeader("Authorization"));
		if authheader == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ok": false,
				"error": "Missing authorization header!",
			})
			return;
		}
		// verify if the scheme and tokens are valid (len, scheme, token)
		authParts := strings.SplitN(authheader, " ", 2);

		if len(authParts) != 2 { // valid split length
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ok": false,
				"error": "Invalid token!",
			})
			return;
		}

		scheme := authParts[0];
		token := authParts[1];

		if !strings.EqualFold(scheme, "Bearer") { // valid scheme(Bearer)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ok": false,
				"error": "Invalid token format!",
			})
			return;
		}


		if token == "" { // token length
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"ok": false,
				"error": "Invalid token!",
			})
			return;
		}

		// decode/parse token
		claims, err := auth.ParseToken(jwtSecret, token);
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ok": false,
				"error": err.Error(),
			})
			return;
		}

		c.Set(ctxUserIdKey, claims.ID);
		c.Set(ctxUserRoleKey, claims.Role)
	
		// next
		c.Next();
	}
}


func GetID (ctx *gin.Context) (string, bool) {
	res, ok := ctx.Get(ctxUserIdKey);
	if !ok {
		return "", false
	}

	userID, _ := res.(string)

	return userID, ok
}

func GetRole (ctx *gin.Context) (string, bool) {
	res, ok := ctx.Get(ctxUserRoleKey);
	if !ok {
		return "", false
	}

	role, _ := res.(string)

	return role, ok
}