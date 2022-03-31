package model

import (

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	DB = db
	return DB
}

func FindPost(post_id interface{}) (*Post, error) {
	var post Post
	if err := DB.Where("id = ?", post_id).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func FetchAllPost(hidden interface{}) ([]Post, error) {
	var posts []Post
	if err := DB.Preload("Comments").Where("hidden = ?", hidden).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func FindUser(user_id interface{}) (*User, error) {
	var user User
	if err := DB.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func AddComment(post_id interface{}, comment_req CommentRequest, user_id uint) (*Post, error) {
	var post Post
	if err := DB.Where("id = ?", post_id).First(&post).Error; err != nil {
		return nil, err
	}
	var comment Comment
	comment.CommentRequest = comment_req
	comment.UserID = user_id
	if err := DB.Model(&post).Association("Comments").Append(&comment); err != nil {
		return nil, err
	}
	
	return &post, nil
}

func SearchPost(keyword string) ([]Post, error) {
	var posts []Post
	if err := DB.Preload("Comments").Where("content LIKE ?", "%"+keyword+"%").Find(&posts).Error;err != nil {
		return nil, err
	}
	return posts, nil
}

func Save(object interface{}) (error) {
	if err := DB.Save(&object).Error;err != nil {
		return err
	}
	return nil
}
