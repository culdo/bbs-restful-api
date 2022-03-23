package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/culdo/bbs-restful-api/model"
)

func HiddenPost(c *gin.Context) {
	postid := c.Param("id")

	var post model.Post
	model.DB.Where("id = ?", postid).First(&post)

	if post.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post id"})
		return
	}

	var hidden_post model.HiddenPost
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hidden_post.Post = post
	model.DB.Save(&hidden_post)
	c.JSON(http.StatusCreated, gin.H{"message": "HiddenPost created successfully!", "HiddenPost": hidden_post})
}

func BanUser(c *gin.Context) {
	userid := c.Param("id")

	var user model.User
	model.DB.Where("id = ?", userid).First(&user)

	if user.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	var banned_user model.BannedUser
	if err := c.ShouldBindJSON(&banned_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	banned_user.User = user
	model.DB.Save(&banned_user)
	c.JSON(http.StatusCreated, gin.H{"message": "BannedUser created successfully!", "BannedUser": banned_user})
}