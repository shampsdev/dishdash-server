package ws

import (
	"dishdash.ru/internal/gateways/ws/newroom"
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

func setupRouter(s *Server, cases usecase.Cases) {
	// room.SetupHandlers(s.WsServer, cases)
	newroom.SetupHandlers(s.WsServer, cases)
	s.Router.GET("/socket.io/*any", gin.WrapH(s.WsServer))
	s.Router.POST("/socket.io/*any", gin.WrapH(s.WsServer))
}
