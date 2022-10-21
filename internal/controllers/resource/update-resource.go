package resource

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type resourceUpdateRequest struct {
	OldWeek   uint32 `json:"oldWeek" binding:"required"`
	NewWeek   uint32 `json:"newWeek" binding:"required"`
	Value     uint32 `json:"value" binding:"required"`
	Email     string `json:"email" binding:"required"`
	GroupName string `json:"name" binding:"required"`
}

func UpdateResource(c *gin.Context) {

	var body resourceUpdateRequest
	var err error

	if err = c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}

	// retrieve user ID and group ID
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

	// Check that the user is in the group
	var users []db.User
	var userID string = ""

	users = group.Users

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

	groupID := group.ID

	// check if the resource entry already exists
	var resource db.Resource

	if body.NewWeek != body.OldWeek {
		res := db.Connection.Where("group_id = ? AND user_id = ? AND week = ?", groupID, userID, body.NewWeek).First(&resource)
		if res.Error == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cannot update, would create duplicate: Entry for requested new week already exists.",
				"data":  nil,
			})
			return
		}
	}

	res := db.Connection.Where("group_id = ? AND user_id = ? AND week = ?", groupID, userID, body.OldWeek).First(&resource)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot update: Entry doesn't exist for this user, in this group, in this week.",
			"data":  nil,
		})
		return
	}

	result := db.Connection.Model(&resource).Update("Week", body.NewWeek).Update("Value", body.Value)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
			"data":  nil,
		})
		return
	}

	// result = db.Connection.Model(&resource).Update("Value", body.Value)
	// if result.Error != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": result.Error.Error(),
	// 		"data":  nil,
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "",
	})
}
