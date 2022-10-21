package group

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)

type deleteGroupRequest struct {
	GroupName string `json:"name" binding:"required"`
}

func DeleteGroup(c *gin.Context) {

	var body deleteGroupRequest
	var err error

	if err = c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}

	// Check if the group exists
	var group db.Group
	res := db.Connection.Preload("Users").Preload("Resource").First(&group, "name=?", body.GroupName)
	if res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Can't Delete Group: group doesn't exist.",
			"data":  nil,
		})
		return
	}

	// remove the group ID from each user
	// var users []db.User
	users := group.Users

	for _, user := range users {
		result := db.Connection.Model(&user).Update("GroupID", "")
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
				"data":  nil,
			})
			return
		}
	}

	// delete each resource record for the group
	resources := group.Resource
	for _, resource := range resources {
		res := db.Connection.Delete(&resource)
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": res.Error.Error(),
				"data":  nil,
			})
			return
		}
	}

	// delete the group record itself
	result := db.Connection.Delete(&group)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
			"data":  nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "success",
	})
}
