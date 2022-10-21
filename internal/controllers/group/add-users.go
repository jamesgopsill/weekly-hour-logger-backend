package group

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)

type addUsersGroupRequest struct {
	GroupName string   `json:"name" binding:"required"`
	Emails    []string `json:"emails" binding:"required"`
}

// function to see if an element is contained within an array
func contains(list []string, input string) bool {
	for _, element := range list {
		if element == input {
			return true
		}
	}
	return false
}

func AddUsers(c *gin.Context) {

	var body addUsersGroupRequest
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

	// check that the users aren't already in the group
	// extract the users from the group
	var users []db.User
	users = group.Users
	// for each user in the list that are in the group
	for _, user := range users {
		// extract their email
		email := user.Email
		// see if their email is in the 'add these people' list
		res := contains(body.Emails, email)
		// if it is, they're already in the group
		if res {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "Already assigned to this group",
				"data":  nil,
			})
			return
		}
	}

	// Check that all the requested users exist, and add them if they do
	for _, email := range body.Emails {
		// Check if the user exists
		var user db.User
		res := db.Connection.Preload("Group").First(&user, "email=?", email)
		if res.Error != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "User does not exist",
				"data":  nil,
			})
			return
		}

		// Check if the user is in another group
		if user.Group.Name != "" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "User is already in a different group",
				"data":  nil,
			})
			return
		}

		// if they do exist and aren't in another group, add them to the user list for the group
		users = append(users, user)
	}

	// update the database with the new users list
	db.Connection.Model(&group).Update("Users", users)

	// pass success
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "success",
	})
}
