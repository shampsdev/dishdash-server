package ws

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"

	"dishdash.ru/internal/usecase"

	socketio "github.com/googollee/go-socket.io"

	"github.com/gin-gonic/gin"
	"github.com/tj/go-spin"
	"golang.org/x/sync/errgroup"
)

const shutdownDuration = 1500 * time.Millisecond

type Server struct {
	Router   *gin.Engine
	WsServer *socketio.Server
}

func NewServer(useCases usecase.Cases, router *gin.Engine) *Server {
	s := &Server{
		Router:   router,
		WsServer: newSocketIOServer(),
	}

	setupRouter(s, useCases)

	return s
}

func (s *Server) Run(ctx context.Context) error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		return s.WsServer.Serve()
	})

	<-ctx.Done()
	err := errors.Join(s.WsServer.Close())
	err = errors.Join(eg.Wait(), err)
	shutdownWait()
	return err
}

func newSocketIOServer() *socketio.Server {
	wt := websocket.Default
	// TODO legal CheckOrigin
	wt.CheckOrigin = func(_ *http.Request) bool {
		return true
	}

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			wt,
		},
	})

	return server
}

func shutdownWait() {
	spinner := spin.New()
	const spinIterations = 20
	for range spinIterations {
		fmt.Printf("\rgraceful shutdown %s ", spinner.Next())
		time.Sleep(shutdownDuration / spinIterations)
	}
	fmt.Println()
}
