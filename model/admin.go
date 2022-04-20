package model

import (
	"log"

	"github.com/culdo/bbs-restful-api/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateAdmin() {
	_, err := FindUserByName("admin")
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err.Error())
	}
	if err == nil {
		log.Println("use previous created Admin")
		return
	}
	adminPass, err := bcrypt.GenerateFromPassword([]byte(config.AdminPasswd), bcrypt.DefaultCost)
	if err != nil {
		panic(err.Error())
	}
	err = DB.Create(&User{Username: "admin", HashedPassword: adminPass}).Error
	if err != nil {
		panic(err.Error())
	}
	log.Println("admin created")
}

func HidePost(pid interface{}, hidden bool) error {
	var post Post
	if err := DB.Where("id = ?", pid).First(&post).Error; err != nil {
		return err
	}
	post.Hidden = hidden
	if err := DB.Save(&post).Error; err != nil {
		return err
	}
	return nil
}

func ActivateUser(uid interface{}, active bool) error {
	var user User
	if err := DB.Where("id = ?", uid).First(&user).Error; err != nil {
		return err
	}
	user.Active = active
	if err := DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}