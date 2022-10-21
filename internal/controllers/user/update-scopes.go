package user

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/utils"
)

type updateScopesRequest struct {
	Email  string   `json:"email" binding:"required"`
	Scopes []string `json:"scopes" binding:"required"`
}

func UpdateScopes(c *gin.Context) {

	var body updateScopesRequest
	var err error

	if err = c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}

	claims, ok := c.MustGet(gin.AuthUserKey).(*MyCustomClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Auth pass-through problem.",
			"data":  nil,
		})
		return
	}

	if !utils.Contains(claims.Scopes, db.ADMIN_SCOPE) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Need adminstrative rights.",
			"data":  nil,
		})
		return
	}

	var user db.User

	res := db.Connection.First(&user, "email=?", body.Email)
	if res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Account does not exist",
			"data":  nil,
		})
		return
	}

	// Reset and update scopes according to request
	user.Scopes = nil
	user.Scopes = append(user.Scopes, db.USER_SCOPE)

	if utils.Contains(body.Scopes, db.ADMIN_SCOPE) {
		user.Scopes = append(user.Scopes, db.ADMIN_SCOPE)
	}

	result := db.Connection.Model(&user).Update("Scopes", user.Scopes)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error: Cannot update scopes",
			"data":  nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "success",
	})

}
