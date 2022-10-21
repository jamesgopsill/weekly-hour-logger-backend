package group

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)

type listUsersGroupRequest struct {
	GroupName string `json:"name" binding:"required"`
}

func ListUsersInGroup(c *gin.Context) {

	var body listUsersGroupRequest
	var err error

	if err = c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}

	// Retrieve the group
	var group db.Group
	groupResult := db.Connection.Preload("Users").First(&group, "name=?", body.GroupName)

	// check that the group exists
	if groupResult.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Group doesn't exist",
			"data":  nil,
		})
		return
	}

	users := group.Users

	// pass success
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  users,
	})
}
