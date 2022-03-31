package model

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(userReq UserRequest) error {
	var user User
	err := DB.Where("username = ?", userReq.Username).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if user.ID > 0 {
		return errors.New("User already exists")
	}
	var newUser User
	newUser.Username = userReq.Username
	newUser.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil{
		return err
	}
	if err := Save(&newUser); err != nil{
		return err
	}
	log.Println(newUser)
	return nil
}