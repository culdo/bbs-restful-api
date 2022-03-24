package controller

import (
	"net/http"

	jwtapple2 "github.com/appleboy/gin-jwt/v2"
	"github.com/culdo/bbs-restful-api/config"
	"github.com/culdo/bbs-restful-api/model"
	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	claims := jwtapple2.ExtractClaims(c)

	var user model.User
	model.DB.Where("id = ?", claims[config.IdentityKey]).First(&user)

	if user.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	var post_req model.PostRequest
	if err := c.ShouldBindJSON(&post_req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post model.Post
	post.PostRequest = post_req
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

	var comment_req model.CommentRequest
	if err := c.ShouldBindJSON(&comment_req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var post model.Post
	model.DB.Where("id = ?", post_id).First(&post)
	
	var comment model.Comment
	comment.CommentRequest = comment_req
	comment.UserID = user.ID
	model.DB.Model(&post).Association("Comments").Append(&comment)
	c.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully!", "Post": post})
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
	var comments []model.Comment
	model.DB.Find(&posts)

	if  _, exists := c.Get("hidden_post"); exists{
		var buff_posts []model.Post
		for _, post := range posts {
			if !post.Hidden {
				buff_posts = append(buff_posts, post)	
			}
		}
		posts = buff_posts
	}

	for i, post := range posts {
		model.DB.Model(&post).Association("Comments").Find(&comments)
		posts[i].Comments = comments
	}

	if len(posts) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No Posts found", "data": posts})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}

