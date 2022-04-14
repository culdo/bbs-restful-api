package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/culdo/bbs-restful-api/model"
	"github.com/culdo/bbs-restful-api/migration"
	"github.com/culdo/bbs-restful-api/router"
)

func init() {
	db := model.Init()
	migration.Migrate(db)
	model.CreateAdmin()
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
