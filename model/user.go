package model

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(userReq UserRequest) error {
	userCheck, err := FindUserByName(userReq.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if userCheck != nil {
		return errors.New("User already exists")
	}
	var user User
	user.Username = userReq.Username
	user.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil{
		return err
	}
	if err := Save(&user); err != nil{
		return err
	}
	log.Println(user)
	return nil
}