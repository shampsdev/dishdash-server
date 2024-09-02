package place

import (
	"dishdash.ru/internal/gateways/http/tag"
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	placeGroup := r.Group("places")
	placeGroup.POST("", SavePlace(cases.Place))
	placeGroup.GET("", GetAllPlaces(cases.Place))
	placeGroup.POST("tags", tag.CreateTag(cases.Tag))
	placeGroup.GET("tags", tag.GetAllTags(cases.Tag))
	placeGroup.PUT("tag/:id", tag.UpdateTag(cases.Tag))
	placeGroup.DELETE("tag/:id", tag.DeleteTag(cases.Tag))
}
