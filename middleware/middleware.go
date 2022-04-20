package middleware

import (
	"github.com/gin-gonic/gin"
)

func DoHidePost(answer bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("doHidePost", answer) 
		c.Next()
	}
}
