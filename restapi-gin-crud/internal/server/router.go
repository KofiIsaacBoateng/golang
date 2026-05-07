package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default();

	r.GET("/hello", func (c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
			"msg": "First route with gin",
		})
	})

	return r
}