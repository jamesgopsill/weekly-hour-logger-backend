package group

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)

func ListGroups(c *gin.Context) {

	// Retrieve the group
	var groups []db.Group
	res := db.Connection.Preload("Users").Find(&groups)

	if res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "No groups found",
			"data":  nil,
		})
		return
	}

	// pass success
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  groups,
	})
}
