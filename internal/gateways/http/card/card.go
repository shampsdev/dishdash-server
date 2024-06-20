package card

import (
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	cardGroup := r.Group("cards")
	cardGroup.POST("", CreateCard(cases.Card))
	cardGroup.GET("", GetAllCards(cases.Card))
	cardGroup.POST("tags", CreateTag(cases.Tag))
	cardGroup.GET("tags", GetAllTags(cases.Tag))
}
