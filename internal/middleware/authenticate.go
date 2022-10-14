package middleware

import (
	"jamesgopsill/resource-logger-backend/internal/config"
	"jamesgopsill/resource-logger-backend/internal/controllers/user"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/utils"
)

// headers are the headers we expect in a auth
type headers struct {
	Authorization string `header:"Authorization"`
}

func Authenticate(scope string) gin.HandlerFunc {

	fn := func(c *gin.Context) {

		var h headers
		var err error
		if err = c.ShouldBindHeader(&h); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Headers issue",
				"data":  nil,
			})
			return
		}

		// Split the authorization on white space. Expecting Bearer ...
		els := strings.Split(h.Authorization, " ")

		// Check its length
		if len(els) != 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid token (1)",
				"data":  nil,
			})
			return
		}

		token, err := jwt.ParseWithClaims(els[1], &user.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			/*
				if _, ok := token.Method.(jwt.SigningMethodHS256); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
			*/
			return config.MySigningKey, nil
		})
		if err != nil {
			log.Info().Err(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid token (2)",
				"data":  nil,
			})
			return
		}

		claims, ok := token.Claims.(*user.MyCustomClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid token (3)",
				"data":  nil,
			})
		}

		if !utils.Contains(claims.Scopes, scope) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Not Authorised",
				"data":  nil,
			})
		}

		// Pass on through the
		c.Set(gin.AuthUserKey, claims)
	}
	return gin.HandlerFunc(fn)
}
