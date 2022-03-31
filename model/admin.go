package model

import "golang.org/x/crypto/bcrypt"

func CreateAdmin() {
	adminPass, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		panic(err.Error())
	}
	err = DB.Create(&User{Username: "admin", HashedPassword: adminPass}).Error
	if err != nil {
		panic(err.Error())
	}
}