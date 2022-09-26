package group

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)

type addResourceGroupRequest struct {
	GroupName string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Week      uint32 `json:"week" binding:"required"`
}

func AddResourceToGroup(c *gin.Context) {

	var body addResourceGroupRequest
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
	groupResult := db.Connection.Preload("Resource").Preload("Users").First(&group, "name=?", body.GroupName)

	// check that the group exists
	if groupResult.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Group doesn't exist",
			"data":  nil,
		})
		return
	}

	// Check that the user is in the group

	var users []db.User
	var groupID string
	var userID string = ""

	users = group.Users
	groupID = group.ID

	// extract the user ID associated with the resource
	for _, user := range users {
		if user.Email == body.Email {
			userID = user.ID
		}
	}

	// check that the user was found in the group
	if userID == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "User not found in this group",
			"data":  nil,
		})
		return
	}

	// find the resource entry that matches group, user ID and week
	var resource db.Resource
	res := db.Connection.Where("groupID = ? AND userID = ? AND Week = ?", groupID, userID, body.Week).Find(&resource)

	// check that it exists
	if res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Could not find resource entry for this user, in this group, in this week. Create the resource first.",
			"data":  nil,
		})
		return
	}

	// check that it hasn't already be added
	for _, res := range group.Resource {
		// if the resource ID matches one already in the group, kick it out
		if resource.ID == res.ID {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "Resource already submitted",
				"data":  nil,
			})
			return
		}

		// add the resource to the group
		db.Connection.Model(&group).Update("Resource", resource)
	}

	// pass success
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "success",
	})
}
