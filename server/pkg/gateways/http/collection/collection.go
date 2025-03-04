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
	collectionGroup.GET(":id", GetCollectionByID(cases.Collection))


	collectionGroup.GET("/preview", GetAllCollectionsPreview(cases.Collection))
	collectionGroup.GET("/preview/:id", GetCollectionPreviewByID(cases.Collection))
}
