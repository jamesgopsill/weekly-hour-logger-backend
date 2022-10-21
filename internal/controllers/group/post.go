package group

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type groupRequest struct {
	Name   string   `json:"name" binding:"required"`
	Emails []string `json:"emails" binding:"required"`
}

func Post(c *gin.Context) {

	var body groupRequest

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}

	var users []db.User
	for _, email := range body.Emails {
		// Check if the user exists
		var user db.User
		res := db.Connection.First(&user, "email=?", email)
		if res.Error != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "Account does not exist",
				"data":  nil,
			})
			return
		}
		log.Info().Msg("Adding user " + user.ID + " to group.")
		users = append(users, user)
	}

	newGroup := db.Group{
		Name:  body.Name,
		Users: users,
	}

	result := db.Connection.Create(&newGroup)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
			"data":  nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  newGroup.ID,
	})

}
