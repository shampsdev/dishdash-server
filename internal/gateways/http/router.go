package http

import (
	"dishdash.ru/docs"
	cardHandler "dishdash.ru/internal/gateways/http/handlers/card"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter(r *gin.Engine, useCases UseCases) {
	r.HandleMethodNotAllowed = true
	v1 := r.Group("/api/v1")
	{
		cardHandler.SetupHandlers(v1, useCases.Card)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
