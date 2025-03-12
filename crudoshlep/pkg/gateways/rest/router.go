package rest

import (
	"github.com/shampsdev/dishdash-server/crudoshlep/docs"
	"github.com/shampsdev/dishdash-server/crudoshlep/pkg/gateways/rest/event"
	"github.com/shampsdev/dishdash-server/crudoshlep/pkg/usecase"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter(s *Server, useCases usecase.Cases) {
	s.Router.HandleMethodNotAllowed = true
	s.Router.Use(AllowOrigin())

	v1 := s.Router.Group("/api/v1")
	v1.Use(Logger())
	{
		event.SetupHandlers(v1, useCases)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
