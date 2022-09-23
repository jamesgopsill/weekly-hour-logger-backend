package main

import (
	"net/http"

	"jamesgopsill/resource-logger-backend/internal/config"
	"jamesgopsill/resource-logger-backend/internal/controllers/group"
	"jamesgopsill/resource-logger-backend/internal/controllers/resource"
	"jamesgopsill/resource-logger-backend/internal/controllers/user"
	"jamesgopsill/resource-logger-backend/internal/db"
	"jamesgopsill/resource-logger-backend/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msg("Starting App")
	r := initialiseApp("data/test.db", gin.ReleaseMode)
	r.Run("localhost:3000")
}

func initialiseApp(dbPath string, mode string) *gin.Engine {
	config.Initalise()
	db.Initialise()

	gin.SetMode(mode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}

	r.Use(cors.New(config))

	r.GET("/ping", pong)
	r.POST("/user/register", user.Register)
	r.POST("/user/login", user.Login)
	r.POST("/user/update", middleware.Authenticate(db.USER_SCOPE), user.Update)
	r.POST("/user/refresh-token", middleware.Authenticate(db.USER_SCOPE), user.RefreshToken)
	r.POST("/user/update-password", middleware.Authenticate(db.USER_SCOPE), user.UpdatePassword)

	// r.POST("/group", middleware.Authenticate(db.USER_SCOPE), group.Post)
	r.POST("/group/create-group", middleware.Authenticate(db.USER_SCOPE), group.CreateGroup)

	r.POST("/resource", middleware.Authenticate(db.USER_SCOPE), resource.Post)

	return r
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "pong",
	})
}