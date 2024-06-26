package user

import (
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary Update a user
// @Description Update a existing user in the database
// @Tags users
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param user body usecase.UserInputExtended true "User data"
// @Success 200 {object} userOutput "Updated user"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /users [put]
func UpdateUser(userUseCase usecase.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userInput usecase.UserInputExtended
		err := c.BindJSON(&userInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userUseCase.UpdateUser(c, userInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, userToOutput(user))
	}
}
