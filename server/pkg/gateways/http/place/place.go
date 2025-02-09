package place

import (
	"dishdash.ru/pkg/gateways/http/middlewares"
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	placeGroup := r.Group("places")
	placeGroup.Use(middlewares.ApiTokenAuth())

	placeGroup.GET("by_url", GetPlaceByURL(cases.Place))
	placeGroup.GET("", GetAllPlaces(cases.Place))
	placeGroup.POST("", SavePlace(cases.Place))
	placeGroup.PUT("", UpdatePlace(cases.Place))
	placeGroup.DELETE(":id", DeletePlace(cases.Place))
}
