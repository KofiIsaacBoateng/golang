package user

import (
	"jwt-auth/internal/app"

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
	}