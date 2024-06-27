package card

import (
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// CreateCard godoc
// @Summary Create a card
// @Description Create a new card in the database
// @Tags cards
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param card body usecase.CardInput true "Card data"
// @Success 200 {object} cardOutput "Saved card"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /cards [post]
func CreateCard(cardUseCase usecase.Card) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cardInput usecase.CardInput
		err := c.BindJSON(&cardInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		card, err := cardUseCase.CreateCard(c, cardInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, cardToOutput(card))
	}
}
