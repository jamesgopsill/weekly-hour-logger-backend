package resource

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type resourceRequest struct {
	Week    int32  `json:"week" binding:"required"`
	Value   int32  `json:"value" binding:"required"`
	UserID  string `json:"userID" binding:"required"`
	GroupID string `json:"groupID" binding:"required"`
}

func Post(c *gin.Context) {

	var body []resourceRequest

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}

	// TODO
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "",
	})

}
