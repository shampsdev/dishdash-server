package http

import (
	"dishdash.ru/docs"
	"dishdash.ru/internal/gateways/http/card"
	"dishdash.ru/internal/gateways/http/lobby"
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter(s *Server, useCases usecase.Cases) {
	s.router.HandleMethodNotAllowed = true
	s.router.Use(allowOriginMiddleware())

	v1 := s.router.Group("/api/v1")
	{
		card.SetupHandlers(v1, useCases)
		lobby.SetupHandlers(v1, useCases)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	s.router.GET("/socket.io/*any", gin.WrapH(s.wsServer))
	s.router.POST("/socket.io/*any", gin.WrapH(s.wsServer))
}
