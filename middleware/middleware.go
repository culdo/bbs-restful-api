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

		admin, err := model.FindUserByID(claims[config.IdentityKey])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if admin.ID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Admin id"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func IsUserActived() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwtapple2.ExtractClaims(c)

		user, err := model.FindUserByID(claims[config.IdentityKey])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !user.Active {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is banned"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func DoHidePost(answer bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("hiddenPost", !answer) 
		c.Next()
	}
}
