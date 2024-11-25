package place

import (
	"dishdash.ru/internal/gateways/http/middlewares"
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	placeGroup := r.Group("places")

	placeGroup.GET("", GetAllPlaces(cases.Place))

	protectedGroup := placeGroup.Group("")
	protectedGroup.Use(middlewares.ApiTokenAuth())

	protectedGroup.POST("", SavePlace(cases.Place))
	protectedGroup.PUT("", UpdatePlace(cases.Place))
	protectedGroup.DELETE(":id", DeletePlace(cases.Place))
}
