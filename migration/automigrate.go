package migration

import (
	"github.com/culdo/bbs-restful-api/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&model.Post{})
	db.AutoMigrate(&model.PostRequest{})
	db.AutoMigrate(&model.Comment{})
	db.AutoMigrate(&model.CommentRequest{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.UserRequest{})
}
