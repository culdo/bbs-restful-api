package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/culdo/bbs-restful-api/model"
)

func HidePost(c *gin.Context) {
	pid := c.Param("id")
	if err := model.HidePost(pid, true); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Post is hidden!", "pid": pid})
}

func UnhidePost(c *gin.Context) {
	pid := c.Param("id")
	if err := model.HidePost(pid, false); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Post is Unhidden!", "pid": pid})
}

func BanUser(c *gin.Context) {
	uid := c.Param("id")
	if err := model.ActivateUser(uid, false); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User is banned!", "uid": uid})
}

func ActivateUser(c *gin.Context) {
	uid := c.Param("id")
	if err := model.ActivateUser(uid, true); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User is active!", "uid": uid})
}