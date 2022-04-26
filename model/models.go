package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	HashedPassword []byte
	Active bool `gorm:"default:true"`
}

type Post struct {
	gorm.Model
	Title    string `json:"title"`
	Content  string `json:"content"`
	Comments []Comment `json:"comments"`
	Hidden bool `gorm:"default:false" json:"hidden"`
	UserID   uint `json:"userid"`
}

type Comment struct {
	gorm.Model
	Content string `json:"content"`
	UserID      uint `json:"userid"`
	PostID      uint `json:"postid"`
}
