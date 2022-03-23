package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

type BannedUser struct {
	gorm.Model
	User User
}

type Post struct {
	gorm.Model
	Title    string `json:"title"`
	Content  string `json:"content"`
	Comments []Comment `json:"comments"`
	UserID   uint   `json:"userid"`
}

type HiddenPost struct {
	gorm.Model
	Post Post
}

type Comment struct {
	gorm.Model
	Content string `json:"content"`
	UserID      uint   `json:"userid"`
	PostID      uint   `json:"postid"`
}
