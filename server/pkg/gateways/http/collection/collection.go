package collection

import (
	"dishdash.ru/pkg/gateways/http/middlewares"
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	collectionGroup := r.Group("collections")
	collectionGroup.Use(middlewares.ApiTokenAuth())

	collectionGroup.GET("", GetAllCollections(cases.Collection))

	collectionGroup.POST("", SaveCollection(cases.Collection))
	collectionGroup.PUT("", UpdateCollection(cases.Collection))
	collectionGroup.DELETE(":id", DeleteCollection(cases.Collection))
}
