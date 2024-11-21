package http

import (
	"dishdash.ru/docs"
	"dishdash.ru/internal/gateways/http/lobby"
	"dishdash.ru/internal/gateways/http/metric"
	"dishdash.ru/internal/gateways/http/middlewares"
	"dishdash.ru/internal/gateways/http/place"
	"dishdash.ru/internal/gateways/http/tag"
	"dishdash.ru/internal/gateways/http/user"
	"dishdash.ru/internal/usecase"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter(s *Server, useCases usecase.Cases) {
	s.Router.HandleMethodNotAllowed = true
	s.Router.Use(middlewares.AllowOriginMiddleware())

	v1 := s.Router.Group("/api/v1")
	metric.AddBasicMetrics(v1)
	v1.Use(middlewares.Logger())
	{
		place.SetupHandlers(v1, useCases)
		lobby.SetupHandlers(v1, useCases)
		user.SetupHandlers(v1, useCases)
		tag.SetupHandlers(v1, useCases)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	s.MetricRouter.HandleMethodNotAllowed = true
	metricV1 := s.MetricRouter.Group("/api/v1")
	metricV1.Use(middlewares.Logger())
	metric.SetupHandlers(metricV1, useCases)
}
