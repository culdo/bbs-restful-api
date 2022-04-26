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

func UpdatePost(pid interface{}, attrs map[string] interface{}) error {
	if err := DB.First(&Post{}, pid).Updates(attrs).Error; err != nil {
		return err
	}
	return nil
}

func UpdateUser(uid interface{}, attrs map[string] interface{}) error {
	if err := DB.First(&User{}, uid).Updates(attrs).Error; err != nil {
		return err
	}
	return nil
}