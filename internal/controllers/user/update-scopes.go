package user

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/utils"
)

type updateScopesRequest struct {
	ID     string   `json:"id" binding:"required"`
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

	if claims.ID != body.ID && utils.Contains(claims.Scopes, db.ADMIN_SCOPE) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Need adminstrative rights.",
			"data":  nil,
		})
	}

	var user db.User

	res := db.Connection.First(&user, "id=?", body.ID)
	if res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Account does not exist",
			"data":  nil,
		})
		return
	}

	scopes := user.Scopes
	scopes = nil
	scopes = append(scopes, db.USER_SCOPE)

	if utils.Contains(body.Scopes, "admin") {
		scopes = append(scopes, db.ADMIN_SCOPE)
	}

	result := db.Connection.Model(&user).Update("Scopes", scopes)

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
