package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/culdo/bbs-restful-api/config"
	"github.com/culdo/bbs-restful-api/model"
	"github.com/gin-gonic/gin"
	goauth "google.golang.org/api/oauth2/v2"
)

func Session(name string) gin.HandlerFunc {
	return sessions.Sessions(name, cookie.NewStore([]byte(config.SessionStoreSecret)))
}


// AuthRequired is a simple middleware to check the session
func AuthRequired(userType string) gin.HandlerFunc {
	return func (c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get(config.IdentityKey)
	if userId == nil || (userType=="admin" && userId != uint(1)) {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
	}
}

// login is a handler that parses a form and checks for specific data
func Login(c *gin.Context) {
	session := sessions.Default(c)

	var loginVals model.UserRequest
	if err := c.ShouldBind(&loginVals); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	user, err := model.Login(loginVals)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed"})
		return 
	}

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}
	
	if 	user.ID > 1 && !user.Active	{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are banned"})
		return
	}

	session.Set(config.IdentityKey, user.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(config.IdentityKey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(config.IdentityKey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func AutoRegister(c *gin.Context) {
	userInfo, _ := c.Get("user")
	var userReq model.UserRequest
	userReq.Username = userInfo.(goauth.Userinfo).Id
	userReq.Password = ""

	if err := model.Register(userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
