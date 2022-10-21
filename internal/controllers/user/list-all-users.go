package user

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)

func ListAllUsers(c *gin.Context) {

	// Retrieve the group
	var users []db.User
	groupResult := db.Connection.Find(&users)

	// check that users were found
	if groupResult.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Error, no results found",
			"data":  nil,
		})
		return
	}

	// pass success
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  users,
	})
}
