package serverRouter

import (
	"jwt-auth/internal/app"
	"jwt-auth/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(app *app.App) *gin.Engine {

	r := gin.New();

	r.Use(gin.Logger())
	
	r.Use(gin.Recovery())


	// welcome
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
			"message": "Welcome to my go user auth learning api!",
		})
	})

	// user routes
	user.NewRouter(r, app)

	return r
}