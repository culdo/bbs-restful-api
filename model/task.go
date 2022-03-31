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

func FindPost(pid interface{}) (*Post, error) {
	var post Post
	if err := DB.Where("id = ?", pid).First(&post).Error; err != nil {
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

func FindUserByID(uid interface{}) (*User, error) {
	var user User
	if err := DB.Where("id = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByName(name interface{}) (*User, error) {
	var user User
	if err := DB.Where("username = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func FindAdmin(uid interface{}) (*User, error) {
	var admin User
	if err := DB.Where("id = ?", uid).First(&admin, "username = ?", "admin").Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func AddComment(pid interface{}, commentReq CommentRequest, uid uint) (*Post, error) {
	var post Post
	if err := DB.Where("id = ?", pid).First(&post).Error; err != nil {
		return nil, err
	}
	var comment Comment
	comment.CommentRequest = commentReq
	comment.UserID = uid
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
