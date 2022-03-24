package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/culdo/bbs-restful-api/model"
)

func HiddenPost(c *gin.Context) {
	postid := c.Param("id")

	var post model.Post
	model.DB.First(&post, postid)

	if post.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post id"})
		return
	}

	post.Hidden = true
	model.DB.Save(&post)
	c.JSON(http.StatusCreated, gin.H{"message": "Post is hidden!", "HiddenPost": post})
}

func BanUser(c *gin.Context) {
	userid := c.Param("id")

	var user model.User
	model.DB.First(&user, userid)

	if user.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	user.Active = false
	model.DB.Save(&user)
	c.JSON(http.StatusCreated, gin.H{"message": "User is banned!", "BannedUser": user})
}