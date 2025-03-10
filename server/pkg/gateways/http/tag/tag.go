package tag

import (
	"dishdash.ru/pkg/gateways/http/middlewares"
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	placeGroup := r.Group("places").Group("tag")

	placeGroup.GET("", GetAllTags(cases.Tag))

	protectedGroup := placeGroup.Group("")
	protectedGroup.Use(middlewares.ApiTokenAuth())

	protectedGroup.POST("", CreateTag(cases.Tag))
	protectedGroup.PUT("", UpdateTag(cases.Tag))
	protectedGroup.DELETE(":id", DeleteTag(cases.Tag))
}
