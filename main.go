package main

import (
	"flag"
	"net/http"
	"testing"

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

// https://github.com/golang/go/issues/31859
var _ = func() bool {
	testing.Init()
	return true
}()

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msg("Starting App")
	r := initialiseApp("data/test.db", gin.ReleaseMode)
	r.Run(":3000")
}

func initialiseApp(dbPath string, mode string) *gin.Engine {

	// Handle the command line flags
	dev := flag.Bool("dev", false, "")
	flag.Parse()

	// Initialise
	config.Initalise(dev)
	db.Initialise()

	// Create API
	gin.SetMode(mode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}

	r.Use(cors.New(config))

	r.GET("/ping", pong)
	r.GET("/user/list-all-users", user.ListAllUsers)
	r.POST("/user/register", user.Register)
	r.POST("/user/login", user.Login)
	r.POST("/user/update", middleware.Authenticate(db.ADMIN_SCOPE), user.Update)
	r.POST("/user/refresh-token", middleware.Authenticate(db.USER_SCOPE), user.RefreshToken)
	r.POST("/user/update-password", middleware.Authenticate(db.USER_SCOPE), user.UpdatePassword)
	r.POST("/user/update-scopes", middleware.Authenticate(db.ADMIN_SCOPE), user.UpdateScopes)
	r.POST("/user/find-user", middleware.Authenticate(db.USER_SCOPE), user.FindUser)
	r.POST("/user/reset-password", middleware.Authenticate(db.ADMIN_SCOPE), user.ResetPassword)

	// r.POST("/group", middleware.Authenticate(db.USER_SCOPE), group.Post)
	r.POST("/group/create-group", middleware.Authenticate(db.ADMIN_SCOPE), group.CreateGroup)
	r.POST("/group/add-users", middleware.Authenticate(db.ADMIN_SCOPE), group.AddUsers)
	r.POST("/group/remove-users", middleware.Authenticate(db.ADMIN_SCOPE), group.RemoveUsers)
	r.POST("/group/list-users-in-group", middleware.Authenticate(db.USER_SCOPE), group.ListUsersInGroup)
	r.POST("/group/delete-group", middleware.Authenticate(db.ADMIN_SCOPE), group.DeleteGroup)
	r.POST("/group/list-resource-in-group", middleware.Authenticate(db.USER_SCOPE), group.ListResourceInGroup)
	r.POST("/group/list-groups", middleware.Authenticate(db.ADMIN_SCOPE), group.ListGroups)
	r.POST("/group/find-group", middleware.Authenticate(db.USER_SCOPE), group.FindGroup)

	r.POST("/resource/add-resource", middleware.Authenticate(db.USER_SCOPE), resource.AddResource)
	r.POST("/resource/update-resource", middleware.Authenticate(db.USER_SCOPE), resource.UpdateResource)
	r.POST("/resource/delete-resource", middleware.Authenticate(db.USER_SCOPE), resource.DeleteResource)

	return r
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "pong",
	})
}
