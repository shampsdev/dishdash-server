package user

import (
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a user from the database by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param id path string true "User ID"
// @Success 200 {object} domain.User "User data"
// @Failure 400 "Bad Request"
// @Failure 500 "pkg Server Error"
// @Router /users/{id} [get]
func GetUserByID(userUseCase usecase.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, err := userUseCase.GetUserByID(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
