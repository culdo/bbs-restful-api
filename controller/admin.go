package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/culdo/bbs-restful-api/model"
)

var postAttrs = []string{"hidden"}
var userAttrs = []string{"active"}

func checkReq(updateReq map[string]interface{}, attrs []string) map[string]interface{} {
	checkedReq := make(map[string]interface{})
	for k, v := range updateReq {
		for _, attr := range attrs {
			if k == attr {
				checkedReq[k] = v
			}
		} 
	}
	return checkedReq
}

func UpdatePost(c *gin.Context) {
	pid := c.Param("id")
	var updateReq map[string]interface{}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateReq = checkReq(updateReq, postAttrs)
	if err := model.UpdatePost(pid, updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Post is updated!", "pid": pid})
}

func UpdateUser(c *gin.Context) {
	uid := c.Param("id")
	var updateReq map[string]interface{}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateReq = checkReq(updateReq, userAttrs)
	if err := model.UpdateUser(uid, updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User is updated!", "uid": uid})
}