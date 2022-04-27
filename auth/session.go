package auth

import (
	"net/http"
	"time"

	"github.com/antonlindstrom/pgstore"
	"github.com/gin-contrib/sessions"

	"github.com/culdo/bbs-restful-api/config"
	"github.com/culdo/bbs-restful-api/model"
	"github.com/gin-gonic/gin"
)

func Session(name string) gin.HandlerFunc {
	sql, err := model.DB.DB()
	if err != nil {
		panic(err.Error())
	}
	store, err := pgstore.NewPGStoreFromPool(sql, []byte(config.SessionKey))

	if err != nil {
		panic(err.Error())
	}

	store.Cleanup(time.Minute * 5)
	return sessions.Sessions(name, &Sstore{store})
}

type Sstore struct {
	*pgstore.PGStore
}

func (s *Sstore) Options(options sessions.Options) {
	s.PGStore.Options = options.ToGorillaOptions()
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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// login is a handler that parses a form and checks for specific data
func Login(c *gin.Context) {
	session := sessions.Default(c)

	var loginVals LoginRequest
	if err := c.ShouldBind(&loginVals); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	user, err := model.Login(loginVals.Username, loginVals.Password)
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
