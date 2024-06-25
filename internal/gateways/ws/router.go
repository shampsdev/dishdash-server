package ws

import (
	"fmt"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"

	socketio "github.com/googollee/go-socket.io"
)

func SetupRouter(s *Server, useCases usecase.Cases) {
	s.WsServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})

	s.WsServer.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		s.Emit("reply", "have "+msg)
		fmt.Printf(msg)
	})

	s.WsServer.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	s.WsServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	s.Router.GET("/socket.io/*any", gin.WrapH(s.WsServer))
	s.Router.POST("/socket.io/*any", gin.WrapH(s.WsServer))
}
