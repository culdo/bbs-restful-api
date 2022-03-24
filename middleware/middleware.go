package middleware

import (
	"net/http"

	"github.com/culdo/bbs-restful-api/config"
	"github.com/culdo/bbs-restful-api/model"
	"github.com/gin-gonic/gin"
	jwtapple2 "github.com/appleboy/gin-jwt/v2"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwtapple2.ExtractClaims(c)

		var user model.User
		model.DB.Where("id = ?", claims[config.IdentityKey]).First(&user, "username = ?", "admin")

		if user.ID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Admin id"})
			c.Abort()
		}

		c.Next()
	}
}

func IsUserActived() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwtapple2.ExtractClaims(c)

		var user model.User
		model.DB.Where("id = ?", claims[config.IdentityKey]).First(&user)

		if !user.Active {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is banned"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func IsPostHidden() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("hidden_post", true) 
		c.Next()
	}
}
