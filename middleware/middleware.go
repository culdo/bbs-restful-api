package middleware

import (
	"net/http"

	jwtapple2 "github.com/appleboy/gin-jwt/v2"
	"github.com/culdo/bbs-restful-api/config"
	"github.com/gin-gonic/gin"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwtapple2.ExtractClaims(c)

		if uint(claims[config.IdentityKey].(float64)) != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Admin"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func DoHidePost(answer bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("doHidePost", answer) 
		c.Next()
	}
}
