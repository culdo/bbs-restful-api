package config

import "os"

const (
	IdentityKey = "id"
	PageItemNum = 30
)
var ( 
	JWTKey     = os.Getenv("JWT_KEY")
	DatabaseUrl = os.Getenv("DATABASE_URL")
	AdminPasswd = os.Getenv("ADMIN_PASSWD")
)