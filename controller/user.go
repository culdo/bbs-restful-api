package controller

import (
	"log"
	"net/http"

	"github.com/culdo/bbs-restful-api/model"
	"github.com/gin-gonic/gin"
)

func RegisterEndpoint(c *gin.Context) {
	var userReq model.UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userCheck, err := model.FindUserByName(userReq.Username)
	userCheck.UserRequest = userReq
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userCheck.ID > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "User already exists"})
		return
	}
	log.Println(userCheck)
	
	if err := model.Save(&userCheck); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
