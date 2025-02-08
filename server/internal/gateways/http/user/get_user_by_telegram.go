package user

import (
	"net/http"
	"strconv"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// GetUserByTelegram godoc
// @Summary Get user by Telegram
// @Description Get a user from the database by Telegram number
// @Tags users
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param telegram path string true "Telegram number"
// @Success 200 {object} domain.User "User data"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /users/telegram/{telegram} [get]
func GetUserByTelegram(userUseCase usecase.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		telegramStr := c.Param("telegram")

		telegram, err := strconv.ParseInt(telegramStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
			return
		}

		user, err := userUseCase.GetUserByTelegram(c, &telegram)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
