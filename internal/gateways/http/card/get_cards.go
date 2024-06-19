package card

import (
	"net/http"

	"dishdash.ru/pkg/filter"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// GetAllCards godoc
// @Summary Get cards
// @Description Get a list of cards from the database
// @Tags cards
// @Accept  json
// @Produce  json
// @Schemes http https
// @Success 200 {array} cardOutput "List of cards"
// @Failure 500
// @Router /cards [get]
func GetAllCards(cardUseCase *usecase.Card) gin.HandlerFunc {
	return func(c *gin.Context) {
		cards, err := cardUseCase.GetAllCards(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, filter.Map(cards, cardToOutput))
	}
}
