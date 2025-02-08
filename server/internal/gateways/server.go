package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"dishdash.ru/internal/gateways/http"
	"dishdash.ru/internal/gateways/ws"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/tj/go-spin"

	"golang.org/x/sync/errgroup"
)

const shutdownDuration = 1500 * time.Millisecond

type Server struct {
	HttpServer *http.Server
	WsServer   *ws.Server
}

func NewServer(useCases usecase.Cases) *Server {
	r := gin.New()
	r.Use(gin.Recovery())

	s := &Server{
		HttpServer: http.NewServer(useCases, r),
		WsServer:   ws.NewServer(useCases, r),
	}

	return s
}

func (s *Server) Run(ctx context.Context) error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		return s.HttpServer.HttpServer.ListenAndServe()
	})
	eg.Go(func() error {
		return s.HttpServer.MetricHttpServer.ListenAndServe()
	})
	eg.Go(func() error {
		return s.WsServer.WsServer.Serve()
	})

	<-ctx.Done()
	err := s.HttpServer.HttpServer.Shutdown(ctx)
	err = errors.Join(err, s.HttpServer.MetricHttpServer.Shutdown(ctx))
	err = errors.Join(err, s.WsServer.WsServer.Close())
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
