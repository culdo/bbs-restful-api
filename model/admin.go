package model

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateAdmin() {
	err := DB.Where("id = ?", 1).First(&User{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err.Error())
	}
	if err == nil {
		log.Println("use previous created Admin")
		return
	}
	adminPass, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		panic(err.Error())
	}
	err = DB.Create(&User{Username: "admin", HashedPassword: adminPass}).Error
	if err != nil {
		panic(err.Error())
	}
	log.Println("admin created")
}