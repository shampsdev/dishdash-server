package ws

import (
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

func setupRouter(s *Server, _ usecase.Cases) {
	s.Router.GET("/socket.io/*any", gin.WrapH(s.WsServer))
	s.Router.POST("/socket.io/*any", gin.WrapH(s.WsServer))
}
