package user

import (
	"dishdash.ru/internal/domain"
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// SaveUser godoc
// @Summary Save a user
// @Description Save a new user in the database
// @Tags users
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param user body domain.User true "User data"
// @Success 200 {object} domain.User "Saved user"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /users [post]
func SaveUser(userUseCase usecase.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := new(domain.User)
		err := c.BindJSON(&user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err = userUseCase.SaveUser(c, user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
