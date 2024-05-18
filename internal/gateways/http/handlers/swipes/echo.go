package swipes

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func SetupEcho(wsServer *socketio.Server) {
	wsServer.OnConnect("", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	wsServer.OnEvent("", "echo", func(s socketio.Conn, msg string) {
		log.Println("echo:", msg)
		s.Emit("echo", msg)
	})

	wsServer.OnError("", func(_ socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	wsServer.OnDisconnect("", func(_ socketio.Conn, msg string) {
		log.Println("closed", msg)
	})
}
