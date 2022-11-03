package user

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)

type findUserRequest struct {
	Email string `json:"email" binding:"required"`
}

func FindUser(c *gin.Context) {

	// Retrieve the group
	var body findUserRequest
	var err error

	if err = c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}

	var user db.User
	res := db.Connection.First(&user, "email", body.Email)

	// check that the user was found
	if res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Error, no user found with that email",
			"data":  nil,
		})
		return
	}

	// pass success
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  user,
	})
}
