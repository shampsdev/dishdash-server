package collection

import (
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	collectionGroup := r.Group("collections")

	collectionGroup.GET("", GetAllCollections(cases.Collection))
	collectionGroup.GET(":id", GetCollectionByID(cases.Collection))

	collectionGroup.GET("/preview", GetAllCollectionsPreview(cases.Collection))
	collectionGroup.GET("/preview/:id", GetCollectionPreviewByID(cases.Collection))
}
