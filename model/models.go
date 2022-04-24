package model

import (
	"gorm.io/gorm"
)
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	HashedPassword []byte
	Active bool `gorm:"default:true"`
}

type PostRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type Post struct {
	gorm.Model
	PostRequest
	Comments []Comment
	Hidden bool `gorm:"default:false"`
	UserID   uint
}

type CommentRequest struct {
	Content string `json:"content"`
}

type Comment struct {
	gorm.Model
	CommentRequest
	UserID      uint
	PostID      uint
}
