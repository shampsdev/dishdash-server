package user

import (
	"net/http"

	"dishdash.ru/internal/domain"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// SaveUserWithID godoc
// @Summary Save a user with specific id
// @Description Save a new user in the database with specific id
// @Tags users
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param user body domain.User true "User data"
// @Success 200 {object} domain.User "Saved user"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /users/with_id [post]
func SaveUserWithID(userUseCase usecase.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := new(domain.User)
		err := c.BindJSON(&user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = userUseCase.SaveUserWithID(c, user, user.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
