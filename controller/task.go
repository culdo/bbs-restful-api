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

	user, err := model.FindUserByID(claims[config.IdentityKey])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.ID <= 0{
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
	if err := model.Save(&post); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully!", "Post": post})
}

func CreateComment(c *gin.Context) {
	post_id := c.Param("id")
	claims := jwtapple2.ExtractClaims(c)
	user, err := model.FindUserByID(claims[config.IdentityKey])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	var comment_req model.CommentRequest
	if err := c.ShouldBindJSON(&comment_req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := model.AddComment(post_id, comment_req, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully!", "Post": post})
}

func FetchAllPost(c *gin.Context) {

	hiddenPost, _ := c.Get("hiddenPost")
	posts, err := model.FetchAllPost(hiddenPost)
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if len(posts) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No Posts found", "data": posts})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func SearchAllPost(c *gin.Context) {
	
	keyword := c.Query("keyword")
	posts, err := model.SearchPost(keyword)
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if len(posts) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No Posts found", "data": posts})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}
