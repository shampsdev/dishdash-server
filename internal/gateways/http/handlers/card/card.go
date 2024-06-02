package card

import (
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cardUseCase *usecase.Card) {
	cardGroup := r.Group("cards")
	cardGroup.GET("", GetCards(cardUseCase))
	cardGroup.POST("", SaveCard(cardUseCase))
}
