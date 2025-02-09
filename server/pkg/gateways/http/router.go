package http

import (
	"dishdash.ru/docs"
	"dishdash.ru/pkg/gateways/http/lobby"
	"dishdash.ru/pkg/gateways/http/metric"
	"dishdash.ru/pkg/gateways/http/middlewares"
	"dishdash.ru/pkg/gateways/http/place"
	"dishdash.ru/pkg/gateways/http/tag"
	"dishdash.ru/pkg/gateways/http/user"
	"dishdash.ru/pkg/usecase"
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
	metricV1 := s.MetricRouter.Group("")
	metricV1.Use(middlewares.Logger())
	metric.SetupHandlers(metricV1, useCases)
}
