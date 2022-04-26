package model

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(username, password string) error {
	userCheck, err := FindUserByName(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if userCheck != nil {
		return errors.New("User already exists")
	}
	var user User
	user.Username = username
	user.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return err
	}
	if err := Save(&user); err != nil{
		return err
	}
	return nil
}