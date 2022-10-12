package main

import (
	"jamesgopsill/resource-logger-backend/internal/controllers/user"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var r *gin.Engine
var invalidSignedString string
var validUserSignedString string
var validUserClaims *user.MyCustomClaims

const SECRET = "test"

type apiResponse struct {
	Error string
	Data  string
}

func init() {
	dbPath := "data/test.db"
	issuer := "www.test.com"
	if _, err := os.Stat(dbPath); err == nil {
		err := os.Remove(dbPath)
		if err != nil {
			panic("Error")
		}
	}
	os.Setenv("GO_REST_JWT_SECRET", SECRET)
	os.Setenv("GO_REST_DB_PATH", dbPath)
	os.Setenv("GO_REST_JWT_ISSUER", issuer)

	var invalidScopes []string
	invalidClaims := user.MyCustomClaims{
		Name:   "a",
		Email:  "b",
		Scopes: invalidScopes,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 24*60*60,
			Issuer:    issuer,
		},
	}

	invalidToken := jwt.NewWithClaims(jwt.SigningMethodHS256, invalidClaims)
	invalidSignedString, _ = invalidToken.SignedString(SECRET)

	r = initialiseApp(dbPath, gin.ReleaseMode)
}
