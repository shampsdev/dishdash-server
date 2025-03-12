package http

import (
	"dashboard.dishdash.ru/docs"
	"dashboard.dishdash.ru/pkg/gateways/http/collection"
	"dashboard.dishdash.ru/pkg/gateways/http/image"
	"dashboard.dishdash.ru/pkg/gateways/http/middlewares"
	"dashboard.dishdash.ru/pkg/gateways/http/place"
	"dashboard.dishdash.ru/pkg/gateways/http/tag"
	"dishdash.ru/pkg/usecase"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter(s *Server, useCases usecase.Cases) {
	s.Router.HandleMethodNotAllowed = true
	s.Router.Use(middlewares.AllowOriginMiddleware())

	v1 := s.Router.Group("/api/v1")
	v1.Use(middlewares.Logger())
	{
		place.SetupHandlers(v1, useCases)
		tag.SetupHandlers(v1, useCases)
		collection.SetupHandlers(v1, useCases)
		image.SetupHandlers(v1)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
