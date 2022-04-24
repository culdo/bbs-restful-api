package jwt

import (
	"net/http"
	"strings"
	"time"

	jwtapple2 "github.com/appleboy/gin-jwt/v2"
	"github.com/culdo/bbs-restful-api/config"
	"github.com/culdo/bbs-restful-api/model"
	"github.com/gin-gonic/gin"
)

func SetupAuth() (*jwtapple2.GinJWTMiddleware, error) {
	authMiddleware, err := jwtapple2.New(&jwtapple2.GinJWTMiddleware{
		Realm:           "bbsapigo",
		Key:             []byte(config.JWTKey),
		Timeout:         time.Hour * 24,
		MaxRefresh:      time.Hour,
		IdentityKey:     config.IdentityKey,
		PayloadFunc:     payload,
		IdentityHandler: identityHandler,
		Authenticator:   authenticator,
		Authorizator:    authorizator,
		Unauthorized:    unauthorized,
		LoginResponse:   loginResponse,
		TokenLookup:     "header: Authorization, query: token, cookie: token",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
		SendCookie:       true,
		CookieHTTPOnly:   true,
		CookieName:       "token",
		CookieSameSite:   http.SameSiteDefaultMode,
	})

	return authMiddleware, err
}

func payload(data interface{}) jwtapple2.MapClaims {
	if v, ok := data.(*model.User); ok {
		return jwtapple2.MapClaims{
			config.IdentityKey: v.ID,
		}
	}

	return jwtapple2.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwtapple2.ExtractClaims(c)
	user, err := model.FindUserByID(claims[config.IdentityKey])
	if err!=nil {
		return false
	}

	return *user
}

func authenticator(c *gin.Context) (interface{}, error) {
	var loginVals model.UserRequest
	if err := c.ShouldBind(&loginVals); err != nil {
		return nil, jwtapple2.ErrMissingLoginValues
	}

	user, err := model.Login(loginVals)
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, jwtapple2.ErrFailedAuthentication
	}
	
	if 	user.ID > 1 && !user.Active	{
		return nil, jwtapple2.ErrForbidden
	}

	return user, nil
}

func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(model.User); ok && v.ID > 0 {
		if strings.HasPrefix(c.FullPath(), "/admin") && v.Username != "admin" {
			return false
		}
		return true
	}

	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"message": message,
	})
}

func loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(code, gin.H{
		"expire": expire,
		"token":  token,
	})
}
