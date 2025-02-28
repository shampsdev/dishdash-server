package collection

import (
	"dashboard.dishdash.ru/cmd/config"
	"dashboard.dishdash.ru/pkg/gateways/http/middlewares"
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	collectionGroup := r.Group("collection")
	collectionGroup.Use(middlewares.ApiTokenAuth(config.C.Auth.ApiToken))

	collectionGroup.GET("", GetAllCollections(cases.Collection))


	collectionGroup.POST("", SaveCollection(cases.Collection))
	collectionGroup.PUT("", UpdateCollection(cases.Collection))
	collectionGroup.DELETE(":id", DeleteCollection(cases.Collection))
}
