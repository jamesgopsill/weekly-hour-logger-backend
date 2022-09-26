package resource

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type resourceRequest struct {
	Week      uint32 `json:"week" binding:"required"`
	Value     uint32 `json:"value" binding:"required"`
	Email     string `json:"email" binding:"required"`
	GroupName string `json:"name" binding:"required"`
}

func AddResource(c *gin.Context) {

	var body resourceRequest
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

	// create a new resource object
	newResource := db.Resource{
		Value:   body.Value,
		Week:    body.Week,
		GroupID: groupID,
		UserID:  userID,
	}

	// add the new resource
	result := db.Connection.Create(&newResource)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
			"data":  nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "",
	})

}
