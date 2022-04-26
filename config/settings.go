package config

import "os"

const (
	IdentityKey = "userid"
	PageItemNum = 30
)
var ( 
	JWTKey     = os.Getenv("JWT_KEY")
	
	DatabaseUrl = os.Getenv("DATABASE_URL")
	AdminPasswd = os.Getenv("ADMIN_PASSWD")
	SessionKey = os.Getenv("SESSION_KEY")

	OAuthClientID = os.Getenv("OAUTH_CLIENT_ID")
	OAuthClientSecret = os.Getenv("OAUTH_CLIENT_SECRET")
)