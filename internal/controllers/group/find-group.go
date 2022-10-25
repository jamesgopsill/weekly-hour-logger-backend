package group

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)

type findGroupRequest struct {
	ID string `json:"id" binding:"required"`
}

func FindGroup(c *gin.Context) {

	// Retrieve the group
	var body findGroupRequest
	var err error

	if err = c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}

	var group db.Group
	res := db.Connection.Preload("Users").First(&group, "id", body.ID)

	// check that the group was found
	if res.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Error, no group found with that ID",
			"data":  nil,
		})
		return
	}

	// pass success
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  group,
	})
}
