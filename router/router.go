package route

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/culdo/bbs-restful-api/auth"
	"github.com/culdo/bbs-restful-api/controller"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	authMiddleware, err := auth.SetupAuth()

	if err != nil {
		log.Fatal("JWT Error" + err.Error())
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to my BBS App"})
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/login", authMiddleware.LoginHandler)

		v1.POST("/register", controller.RegisterEndpoint)

		bbs := v1.Group("/bbs")
		{
			bbs.POST("/post", authMiddleware.MiddlewareFunc(), controller.CreatePost)
			bbs.POST("/post/:id/comment", authMiddleware.MiddlewareFunc(), controller.CreateComment)
			bbs.GET("/post", authMiddleware.MiddlewareFunc(), controller.FetchAllPost)
			// bbs.PUT("/:id", authMiddleware.MiddlewareFunc(), controller.UpdateTask)
		}

		admin := v1.Group("/admin")
		{
			admin.POST("/post/:id/hidden", authMiddleware.MiddlewareFunc(), controller.HiddenPost)
			admin.POST("/user/:id/ban", authMiddleware.MiddlewareFunc(), controller.BanUser)
		}
	}

	authorization := router.Group("/auth")
	authorization.GET("/refresh_token", authMiddleware.RefreshHandler)

	return router
}
