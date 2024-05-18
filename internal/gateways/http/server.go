package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tj/go-spin"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

const shutdownDuration = 1500 * time.Millisecond

type Server struct {
	http.Server
}

func NewServer(options ...func(*Server)) *Server {
	r := gin.Default()
	setupRouter(r)

	s := &Server{
		Server: http.Server{
			Addr:    fmt.Sprintf(":%d", 8080),
			Handler: r,
		},
	}
	for _, o := range options {
		o(s)
	}

	return s
}

func WithPort(port uint16) func(*Server) {
	return func(s *Server) {
		s.Addr = fmt.Sprintf(":%d", port)
	}
}

func (s *Server) Run(ctx context.Context) error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		return s.ListenAndServe()
	})

	<-ctx.Done()
	err := s.Shutdown(ctx)
	err = errors.Join(eg.Wait(), err)
	shutdownWait()
	return err
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
