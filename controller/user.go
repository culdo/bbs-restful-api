package controller

import (
	"log"
	"net/http"

	"github.com/culdo/bbs-restful-api/model"
	"github.com/gin-gonic/gin"
)

func RegisterEndpoint(c *gin.Context) {
	var user_req model.UserRequest
	if err := c.ShouldBindJSON(&user_req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userCheck model.User
	userCheck.UserRequest = user_req
	log.Println(userCheck)
	model.DB.First(&userCheck, "username = ?", userCheck.Username)
	
	if userCheck.ID > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "User already exists"})
		return
	}
	model.DB.Save(&userCheck)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
