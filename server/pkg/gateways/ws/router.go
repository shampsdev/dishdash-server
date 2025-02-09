package ws

import (
	"dishdash.ru/pkg/gateways/ws/room"
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

func setupRouter(s *Server, cases usecase.Cases) {
	room.SetupHandlers(s.WsServer, cases)
	s.Router.GET("/socket.io/*any", gin.WrapH(s.WsServer))
	s.Router.POST("/socket.io/*any", gin.WrapH(s.WsServer))
}
