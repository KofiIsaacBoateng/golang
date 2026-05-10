package user

import (
	"jwt-auth/internal/app"
	"jwt-auth/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, app *app.App) {
	// create repos, services and handlers
	repo := NewRepo(app.MongoDatabase);
	svc := NewService(repo, app.Config.JWT_SECRET);
	handler := NewHandler(svc);

	authGroup := r.Group("/auth");

	{
		authGroup.POST("/register", handler.RegisterUser)
		authGroup.POST("/login", handler.LoginUser)
	}

	opsGroup := r.Group("/user");

	opsGroup.Use(middleware.Authorized(app.Config.JWT_SECRET))

	// eg: get all users routes and update a users profile should be allowed when the user is logged in. 
	opsGroup.GET("/", func(c *gin.Context) {
		ID, _ := middleware.GetID(c)
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
			"message": "Welcome to a protected route meant for only logged in users!",
			"userID": ID,
		})
	})

	// delete operations should only be done by admins
	opsGroup.DELETE("/", middleware.OnlyAdmin(), func(c *gin.Context) {
		role, _ := middleware.GetRole(c)
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
			"message": "Welcome to admin only Route ADMIN!",
			"role": role,
		})
	})
}