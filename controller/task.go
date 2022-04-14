package controller

import (
	"net/http"
	"strconv"

	jwtapple2 "github.com/appleboy/gin-jwt/v2"
	"github.com/culdo/bbs-restful-api/config"
	"github.com/culdo/bbs-restful-api/model"
	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	claims := jwtapple2.ExtractClaims(c)

	uid := uint(claims[config.IdentityKey].(float64))

	var postReq model.PostRequest
	if err := c.ShouldBindJSON(&postReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post model.Post
	post.PostRequest = postReq
	post.UserID = uid
	if err := model.Save(&post); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully!", "Post": post})
}

func CreateComment(c *gin.Context) {
	pid := c.Param("id")
	claims := jwtapple2.ExtractClaims(c)
	uid := uint(claims[config.IdentityKey].(float64))


	var comment_req model.CommentRequest
	if err := c.ShouldBindJSON(&comment_req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := model.CreateComment(pid, comment_req, uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully!", "Post": post})
}

func FetchPosts(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	offset := (page - 1) * config.PageItemNum
	doHidePost, _ := c.Get("doHidePost")
	posts, err := model.FetchPosts(doHidePost.(bool), config.PageItemNum, offset)
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
