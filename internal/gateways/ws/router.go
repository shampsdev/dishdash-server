package ws

import (
	"dishdash.ru/internal/gateways/ws/swipes"
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

func setupRouter(s *Server, useCases usecase.Cases) {
	swipes.SetupHandlers(s.WsServer, useCases)

	s.Router.GET("/socket.io/*any", gin.WrapH(s.WsServer))
	s.Router.POST("/socket.io/*any", gin.WrapH(s.WsServer))
}
