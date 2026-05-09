package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterUser(c *gin.Context) {
	// bind req input to reqbody
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok": false,
			"error": fmt.Sprintf("Invalid req body: %s", err.Error()),
		})
		return;
	}

	// register user
	result, err := h.service.Register(c.Request.Context(), input);
	if(err != nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ok": false,
			"error": err.Error(),
		})
		return;
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"message": "User created successfully!",
		"data": result.u,
		"token": result.Token,
	})


}


func (h *Handler) LoginUser(c *gin.Context) {
	// bind input to request body
	var input LoginInput;

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok": false,
			"error": fmt.Errorf("Invalid JSON format: %w", err),
		})
		return;
	}


	result, err := h.service.Login(c.Request.Context(), input);
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ok": false,
			"error": err.Error(),
		})
		return;
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"message": "Login successful!",
		"data": result.u,
		"token": result.Token,
	})
}