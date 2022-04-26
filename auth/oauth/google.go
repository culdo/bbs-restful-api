// Package google provides you access to Google's OAuth2
// infrastructure. The implementation is based on this blog post:
// http://skarlso.github.io/2016/06/12/google-signin-with-go/
package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/culdo/bbs-restful-api/config"
	"github.com/culdo/bbs-restful-api/model"
	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	goauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Credentials stores google client-ids.
type Credentials struct {
	ClientID     string `json:"clientid"`
	ClientSecret string `json:"secret"`
}

const (
	stateKey  = "state"
	sessionID = "ginoauth_google_session"
)

var (
	conf  *oauth2.Config
)

func init() {
	gob.Register(goauth.Userinfo{})
}

func randToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		glog.Fatalf("[Gin-OAuth] Failed to read rand: %v", err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

// Setup the authorization path
func Setup(redirectURL string, scopes []string) {
	conf = &oauth2.Config{
		ClientID:     config.OAuthClientID,
		ClientSecret: config.OAuthClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
}

func LoginHandler(ctx *gin.Context) {
	stateValue := randToken()
	session := sessions.Default(ctx)
	session.Set(stateKey, stateValue)
	session.Save()
	ctx.Writer.Write([]byte(`
	<html>
		<head>
			<title>Golang Google</title>
		</head>
	  <body>
			<a href='` + GetLoginURL(stateValue) + `'>
				<button>Login with Google!</button>
			</a>
		</body>
	</html>`))
}

func GetLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func AutoRegister(c *gin.Context) {
	userInfo, _ := c.Get("user")
	
	if err := model.Register(userInfo.(goauth.Userinfo).Id, ""); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Handle the exchange code to initiate a transport.
		session := sessions.Default(ctx)

		existingSession := session.Get(sessionID)
		if userInfo, ok := existingSession.(goauth.Userinfo); ok {
			ctx.Set("user", userInfo)
			ctx.Next()
			return
		}

		retrievedState := session.Get(stateKey)
		if retrievedState != ctx.Query(stateKey) {
			ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid session state: %s", retrievedState))
			return
		}

		tok, err := conf.Exchange(context.TODO(), ctx.Query("code"))
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to exchange code for oauth token: %w", err))
			return
		}

		oAuth2Service, err := goauth.NewService(ctx, option.WithTokenSource(conf.TokenSource(ctx, tok)))
		if err != nil {
			glog.Errorf("[Gin-OAuth] Failed to create oauth service: %v", err)
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create oauth service: %w", err))
			return
		}

		userInfo, err := oAuth2Service.Userinfo.Get().Do()
		if err != nil {
			glog.Errorf("[Gin-OAuth] Failed to get userinfo for user: %v", err)
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get userinfo for user: %w", err))
			return
		}

		ctx.Set("user", userInfo)

		session.Set(sessionID, userInfo)
		if err := session.Save(); err != nil {
			glog.Errorf("[Gin-OAuth] Failed to save session: %v", err)
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to save session: %v", err))
			return
		}
	}
}