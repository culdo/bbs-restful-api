package controller

import (
	"net/http"

	"github.com/culdo/bbs-restful-api/model"
	"github.com/gin-gonic/gin"
)
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterEndpoint(c *gin.Context) {
	var userReq RegisterRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := model.Register(userReq.Username, userReq.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
