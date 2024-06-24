package http

import (
	"dishdash.ru/docs"
	server "dishdash.ru/internal/gateways"
	"dishdash.ru/internal/gateways/http/card"
	"dishdash.ru/internal/gateways/http/lobby"
	"dishdash.ru/internal/usecase"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter(s *server.Server, useCases usecase.Cases) {
	s.Router.HandleMethodNotAllowed = true
	s.Router.Use(allowOriginMiddleware())

	v1 := s.Router.Group("/api/v1")
	{
		card.SetupHandlers(v1, useCases)
		lobby.SetupHandlers(v1, useCases)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
