package server

import (
	"net/http"
	"notes-api/internal/notes-api"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewRouter(db *mongo.Database) *gin.Engine {
	r := gin.Default();

	r.GET("/hello", func (c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
			"msg": "First route with gin",
		})
	})

	notes.RegisterRoutes(r, db)

	return r
}