package user

import (
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary Create a user
// @Description Create a new user in the database
// @Tags users
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param user body usecase.UserInput true "User data"
// @Success 200 {object} userOutput "Saved user"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /users [post]
func CreateUser(userUseCase usecase.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userInput usecase.UserInput
		err := c.BindJSON(&userInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userUseCase.CreateUser(c, userInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, userToOutput(user))
	}
}
