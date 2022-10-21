package group

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	// "github.com/google/uuid"
)

type removeUsersGroupRequest struct {
	GroupName string   `json:"name" binding:"required"`
	Emails    []string `json:"emails" binding:"required"`
}

func RemoveUsers(c *gin.Context) {

	var body removeUsersGroupRequest
	var err error

	if err = c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}

	// Check that all the requested users exist
	for _, email := range body.Emails {
		// Check if the user exists
		var user db.User
		res := db.Connection.First(&user, "email=?", email)
		if res.Error != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "User(s) do not exist",
				"data":  nil,
			})
			return
		}
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

	// extract the users from the group
	var users []db.User          // empty slice to add the not-removed users to
	usersOriginal := group.Users // slice of the users previously in the group
	var usersOriginalEmails []string

	// for each user thats previously in the group
	for _, user := range usersOriginal {

		// extract their email
		email := user.Email
		usersOriginalEmails = append(usersOriginalEmails, email)

		// see if their email is in the 'remove these people' list
		res := contains(body.Emails, email)

		// if it is, remove their group ID
		if res {
			// find the user entry for the removed person, and remove their group association
			var removeUser db.User
			db.Connection.First(&removeUser, "email=?", email)
			removeUser.GroupID = ""
			db.Connection.Save(&removeUser)
		} else {
			// if they aren't on the remove list, add them to the new users list
			users = append(users, user)
		}
	}

	// Double check if there were any 'to-remove' people who weren't in the group
	for _, email := range body.Emails {
		res := contains(usersOriginalEmails, email)
		if !res {
			log.Info().Msg("WARNING: User email: " + email + " not in group in the first place.")
		}
	}

	// update the database with the new users list
	db.Connection.Model(&group).Update("Users", users)

	// pass success
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "success",
	})
}
