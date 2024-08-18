package user

import (
	"net/http"

	"dishdash.ru/internal/domain"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

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
