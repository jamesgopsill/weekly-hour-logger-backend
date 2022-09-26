package group

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)

type listResourceGroupRequest struct {
	GroupName string `json:"name" binding:"required"`
}

func ListResourceInGroup(c *gin.Context) {

	var body listResourceGroupRequest
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
	groupResult := db.Connection.Preload("Resource").First(&group, "name=?", body.GroupName)

	// check that the group exists
	if groupResult.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Group doesn't exist",
			"data":  nil,
		})
		return
	}

	resource := group.Resource

	// pass success
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  resource,
	})
}
