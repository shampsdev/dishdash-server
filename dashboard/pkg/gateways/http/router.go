package http

import (
	"dashboard.dishdash.ru/docs"
	"dashboard.dishdash.ru/pkg/gateways/http/middlewares"
	"dashboard.dishdash.ru/pkg/gateways/http/photo"
	"dashboard.dishdash.ru/pkg/gateways/http/place"
	"dashboard.dishdash.ru/pkg/gateways/http/task"
	"dashboard.dishdash.ru/pkg/repo"
	"dishdash.ru/pkg/usecase"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter(s *Server, useCases usecase.Cases, taskRepo repo.Task) {
	s.Router.HandleMethodNotAllowed = true
	s.Router.Use(middlewares.AllowOriginMiddleware())

	v1 := s.Router.Group("/api/v1")
	v1.Use(middlewares.Logger())
	{
		place.SetupHandlers(v1, useCases)
		photo.SetupHandlers(v1)
		task.SetupHandlers(v1, taskRepo)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
