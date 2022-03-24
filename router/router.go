package router

import (
	"log"
	"net/http"

	"github.com/culdo/bbs-restful-api/auth"
	"github.com/culdo/bbs-restful-api/controller"
	"github.com/culdo/bbs-restful-api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	authMiddleware, err := auth.SetupAuth()

	if err != nil {
		log.Fatal("JWT Error" + err.Error())
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to my BBS App"})
	})

	router.POST("/login", authMiddleware.LoginHandler)

	router.POST("/register", controller.RegisterEndpoint)

	bbs := router.Group("/bbs")
	bbs.Use(authMiddleware.MiddlewareFunc())
	bbs.Use(middleware.IsUserActived())
	{
		bbs.POST("/posts", controller.CreatePost)
		bbs.GET("/posts", middleware.IsPostHidden(), controller.FetchAllPost)
		bbs.POST("/posts/:id/comment", controller.CreateComment)
	}

	admin := router.Group("/admin")
	admin.Use(authMiddleware.MiddlewareFunc())
	admin.Use(middleware.IsAdmin())
	{
		admin.GET("/posts", controller.FetchAllPost)
		admin.GET("/posts/:id/hide", controller.HiddenPost)
		admin.GET("/users/:id/ban", controller.BanUser)
	}

	authorization := router.Group("/auth")
	authorization.GET("/refresh_token", authMiddleware.RefreshHandler)

	return router
}
