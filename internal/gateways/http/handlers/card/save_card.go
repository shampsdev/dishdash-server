package card

import (
	"net/http"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/dto"
	"dishdash.ru/internal/filter"
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// SaveCard godoc
// @Summary Save a card
// @Description Save a new card to the database
// @Tags cards
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param card body dto.CardToCreate true "Card data"
// @Success 200 {object} dto.Card "Saved card"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /cards [post]
func SaveCard(cardUseCase *usecase.Card) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cardDto dto.CardToCreate
		err := c.BindJSON(&cardDto)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		errs := dto.ValidateCardToCreate(cardDto)
		if errs != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": filter.Map(errs, func(err validator.FieldError) string {
					return err.Error()
				}),
			})
			return
		}

		var card *domain.Card
		err = card.ParseDtoToCreate(cardDto)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err = cardUseCase.SaveCard(c, card)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusCreated, card.ToDto())
	}
}
