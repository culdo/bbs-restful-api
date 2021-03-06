package controller

import (
	"net/http"
	"strconv"

	"github.com/culdo/bbs-restful-api/config"
	"github.com/culdo/bbs-restful-api/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PostRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type CommentRequest struct {
	Content string `json:"content"`
}

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to my BBS APP"})
}

func CreatePost(c *gin.Context) {
	session := sessions.Default(c)

	uid := session.Get(config.IdentityKey).(uint)

	var postReq PostRequest
	if err := c.ShouldBindJSON(&postReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post model.Post
	post.Title = postReq.Title
	post.Content = postReq.Content

	post.UserID = uid
	if err := model.Save(&post); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully!", "post": post})
}

func CreateComment(c *gin.Context) {
	pid := c.Param("id")
	session := sessions.Default(c)
	uid := session.Get(config.IdentityKey).(uint)


	var comment_req CommentRequest
	if err := c.ShouldBindJSON(&comment_req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var comment model.Comment
	comment.Content = comment_req.Content
	comment.UserID = uid
	_, err := model.CreateComment(pid, comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully!"})
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
