package controller

import (
	"net/http"

	jwtapple2 "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/culdo/bbs-restful-api/config"
	"github.com/culdo/bbs-restful-api/model"
)

func CreatePost(c *gin.Context) {
	claims := jwtapple2.ExtractClaims(c)

	var user model.User
	model.DB.Where("id = ?", claims[config.IdentityKey]).First(&user)

	if user.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	var post model.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.UserID = user.ID
	model.DB.Save(&post)
	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully!", "Post": post})
}

func CreateComment(c *gin.Context) {
	post_id := c.Param("id")
	
	claims := jwtapple2.ExtractClaims(c)

	var user model.User
	model.DB.Where("id = ?", claims[config.IdentityKey]).First(&user)

	if user.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	var post model.Post
	model.DB.Where("id = ?", post_id).First(&post)
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	comment.UserID = user.ID
	model.DB.Model(&post).Association("Comments").Append(comment)
	model.DB.Save(&post)
	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully!", "Post": post})
}

func FetchAllPost(c *gin.Context) {
	claims := jwtapple2.ExtractClaims(c)

	var user model.User
	model.DB.Where("id = ?", claims[config.IdentityKey]).First(&user)

	if user.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	var posts []model.Post
	model.DB.Find(&posts)

	if len(posts) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No Posts found", "data": posts})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}

