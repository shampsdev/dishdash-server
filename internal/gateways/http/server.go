package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"dishdash.ru/cmd/server/config"
	"dishdash.ru/internal/usecase"

	"github.com/tj/go-spin"
	"github.com/gin-gonic/gin"

	"golang.org/x/sync/errgroup"
)

const shutdownDuration = 1500 * time.Millisecond

type Server struct {
	HttpServer http.Server
	Router     *gin.Engine
}

func NewServer(useCases usecase.Cases) *Server {
	r := gin.Default()

	s := &Server{
		Router: r,
		HttpServer: http.Server{
			Addr:    fmt.Sprintf(":%d", config.C.Server.Port),
			Handler: r,
		},
	}

	return s
}

func (s *Server) Run(ctx context.Context) error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		return s.HttpServer.ListenAndServe()
	})

	<-ctx.Done()
	err := s.HttpServer.Shutdown(ctx)
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
