package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"

	"dishdash.ru/internal/usecase"

	socketio "github.com/googollee/go-socket.io"

	"github.com/tj/go-spin"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

const shutdownDuration = 1500 * time.Millisecond

type Server struct {
	httpServer  http.Server
	router      *gin.Engine
	wsServer    *socketio.Server
	allowOrigin string
}

func NewServer(useCases usecase.Cases, options ...func(*Server)) *Server {
	r := gin.Default()

	s := &Server{
		router: r,
		httpServer: http.Server{
			Addr:    fmt.Sprintf(":%d", 8080),
			Handler: r,
		},
		wsServer: newSocketIOServer(),
	}

	for _, o := range options {
		o(s)
	}
	setupRouter(s, useCases)

	return s
}

func WithPort(port uint16) func(*Server) {
	return func(s *Server) {
		s.httpServer.Addr = fmt.Sprintf(":%d", port)
	}
}

func WithAllowOrigin(allowOrigin string) func(*Server) {
	return func(s *Server) {
		s.allowOrigin = allowOrigin
	}
}

func (s *Server) Run(ctx context.Context) error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		return s.httpServer.ListenAndServe()
	})
	eg.Go(func() error {
		return s.wsServer.Serve()
	})

	<-ctx.Done()
	err := s.httpServer.Shutdown(ctx)
	err = errors.Join(eg.Wait(), err)
	shutdownWait()
	return err
}

func newSocketIOServer() *socketio.Server {
	pt := polling.Default

	wt := websocket.Default
	// TODO legal CheckOrigin
	wt.CheckOrigin = func(_ *http.Request) bool {
		return true
	}

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			pt,
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
