package main

import (
	"jamesgopsill/resource-logger-backend/internal/controllers/user"
	"jamesgopsill/resource-logger-backend/internal/db"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

var r *gin.Engine
var invalidSignedString string
var mockAdminSignedString string

var validUserSignedString string
var validUserClaims *user.MyCustomClaims

const SECRET = "test"

type apiResponse struct {
	Error string
	Data  string
}

func init() {
	var err error

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

	// Create a mock invalid JWT token
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
	invalidSignedString, err = invalidToken.SignedString([]byte(SECRET))
	if err != nil {
		log.Error().Err(err).Msg("invalidSignedString Signing Error")
		panic("Init test error")
	}

	// Create a mock admin JWT token
	var adminScopes []string
	adminScopes = append(adminScopes, db.USER_SCOPE, db.ADMIN_SCOPE)
	adminClaims := user.MyCustomClaims{
		Name:   "mock admin",
		Email:  "admin@mock",
		Scopes: adminScopes,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 24*60*60,
			Issuer:    issuer,
		},
	}

	adminToken := jwt.NewWithClaims(jwt.SigningMethodHS256, adminClaims)
	mockAdminSignedString, err = adminToken.SignedString([]byte(SECRET))
	if err != nil {
		log.Error().Err(err).Msg("mockAdminSignedString Signing Error")
		panic("Init test error")
	}

	r = initialiseApp(dbPath, gin.ReleaseMode)
}
