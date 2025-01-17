package user

import (
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// GetAllUsers godoc
// @Summary Get users
// @Description Get a list of users from the database
// @Tags users
// @Accept  json
// @Produce  json
// @Schemes http https
// @Success 200 {array} domain.User "List of users"
// @Failure 500
// @Router /users [get]
// @Security ApiKeyAuth
func GetAllUsers(userUseCase usecase.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := userUseCase.GetAllUsers(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, users)
	}
}
