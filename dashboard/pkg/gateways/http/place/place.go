package place

import (
	"dashboard.dishdash.ru/cmd/config"
	"dashboard.dishdash.ru/pkg/gateways/http/middlewares"
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	placeGroup := r.Group("places")
	placeGroup.Use(middlewares.ApiTokenAuth(config.C.Auth.ApiToken))

	placeGroup.
		GET("by_url", GetPlaceByURL(cases.Place)).
		POST("parse", ParsePlace()).
		GET("", GetAllPlaces(cases.Place)).
		POST("", SavePlace(cases.Place)).
		PUT("", UpdatePlace(cases.Place)).
		GET("id/:id", GetPlaceByID(cases.Place)).
		DELETE("id/:id", DeletePlace(cases.Place)).
		PATCH("", PatchPlace(cases.Place))
}
