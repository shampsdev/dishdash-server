package http

import (
	"dishdash.ru/docs"
	cardHandler "dishdash.ru/internal/gateways/http/handlers/card"
	lobbyHandler "dishdash.ru/internal/gateways/http/handlers/lobby"
	"dishdash.ru/internal/gateways/http/handlers/swipes"
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter(s *Server, useCases usecase.Cases) {
	s.router.HandleMethodNotAllowed = true
	v1 := s.router.Group("/api/v1")
	{
		cardHandler.SetupHandlers(v1, useCases.Card)
		lobbyHandler.SetupHandlers(v1, useCases.Lobby)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	s.router.Use(allowOriginMiddleware(s.allowOrigin))
	s.router.GET("/socket.io/*any", gin.WrapH(s.wsServer))
	s.router.POST("/socket.io/*any", gin.WrapH(s.wsServer))
	swipes.SetupLobby(s.wsServer, useCases)
}
