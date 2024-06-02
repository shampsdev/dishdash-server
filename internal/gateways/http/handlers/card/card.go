package card

import (
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cardUseCase *usecase.Card, tagUseCase *usecase.Tag) {
	cardGroup := r.Group("cards")
	cardGroup.GET("", GetCards(cardUseCase, tagUseCase))
	cardGroup.POST("", SaveCard(cardUseCase))
}
