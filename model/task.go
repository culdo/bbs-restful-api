package model

import (
	"errors"

	"github.com/culdo/bbs-restful-api/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.DatabaseUrl))
	if err != nil {
		panic(err.Error())
	}

	DB = db
	return DB
}

func Login(userReq UserRequest) (*User, error) {
	user, err := FindUserByName(userReq.Username)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(userReq.Password)); err != nil {
		return nil, errors.New("password incorrect")
	}
	return user, nil
}

func FindPostByID(pid interface{}) (*Post, error) {
	var post Post
	if err := DB.Where("id = ?", pid).Find(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func FetchPosts(doHidePost bool, limit int, offset int) ([]Post, error) {
	var posts []Post
	tx := DB.Limit(limit).Offset(offset).Preload("Comments")
	if doHidePost {
		if err := tx.Where("hidden = ?", false).Find(&posts).Error; err != nil {
			return nil, err
		}
	} else {
		if err := tx.Find(&posts).Error; err != nil {
			return nil, err
		}
	}
	return posts, nil
}

func FindUserByID(uid interface{}) (*User, error) {
	var user User
	if err := DB.Where("id = ?", uid).Find(&user).Error; err != nil {
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

func CreateComment(pid interface{}, commentReq CommentRequest, uid uint) (*Post, error) {
	var post Post
	if err := DB.Where("id = ?", pid).Find(&post).Error; err != nil {
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
	if err := DB.Preload("Comments").Where("content LIKE ?", "%"+keyword+"%").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func Save(object interface{}) error {
	if err := DB.Save(object).Error; err != nil {
		return err
	}
	return nil
}
