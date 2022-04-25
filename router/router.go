package router

import (
	"github.com/culdo/bbs-restful-api/auth"
	"github.com/culdo/bbs-restful-api/controller"
	"github.com/culdo/bbs-restful-api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(auth.Session("bbssession"))
	router.GET("/", controller.Index)

	router.POST("/login", auth.Login)
	router.GET("/logout", auth.Logout)
	router.POST("/register", controller.RegisterEndpoint)

	router.GET("/posts", middleware.DoHidePost(true), controller.FetchPosts)
	router.Use(auth.AuthRequired("user")) 
	{
		router.POST("/posts", controller.CreatePost)
		router.POST("/posts/:id/comments", controller.CreateComment)
	}

	admin := router.Group("/admin")
	admin.Use(auth.AuthRequired("admin"))
	{
		admin.GET("/posts", middleware.DoHidePost(false), controller.FetchPosts)
		admin.GET("/posts/search", controller.SearchAllPost)
		admin.POST("/posts/:id/hide", controller.HidePost)
		admin.POST("/posts/:id/unhide", controller.UnhidePost)
		admin.POST("/users/:id/ban", controller.BanUser)
		admin.POST("/users/:id/activate", controller.ActivateUser)
	}

	return router
}
