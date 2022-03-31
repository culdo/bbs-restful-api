package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/culdo/bbs-restful-api/model"
)

func HidePost(c *gin.Context) {
	pid := c.Param("id")

	post, err := model.FindPost(pid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if post.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post id"})
		return
	}

	post.Hidden = true
	model.DB.Save(&post)
	c.JSON(http.StatusCreated, gin.H{"message": "Post is hidden!", "HiddenPost": post})
}

func UnhidePost(c *gin.Context) {
	pid := c.Param("id")

	post, err := model.FindPost(pid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if post.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post id"})
		return
	}

	post.Hidden = false
	model.DB.Save(&post)
	c.JSON(http.StatusCreated, gin.H{"message": "Post is Unhidden!", "UnhiddenPost": post})
}

func BanUser(c *gin.Context) {
	uid := c.Param("id")

	user, err := model.FindUserByID(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	user.Active = false
	model.DB.Save(&user)
	c.JSON(http.StatusCreated, gin.H{"message": "User is banned!", "BannedUser": user})
}

func ActivateUser(c *gin.Context) {
	uid := c.Param("id")

	user, err := model.FindUserByID(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if user.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	user.Active = true
	model.DB.Save(&user)
	c.JSON(http.StatusCreated, gin.H{"message": "User is active!", "ActiveUser": user})
}