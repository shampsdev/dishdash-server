package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"dishdash.ru/cmd/server/config"
	"dishdash.ru/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/tj/go-spin"

	"golang.org/x/sync/errgroup"
)

const shutdownDuration = 1500 * time.Millisecond

type Server struct {
	HttpServer       http.Server
	MetricHttpServer http.Server
	Router           *gin.Engine
	MetricRouter     *gin.Engine
}

func NewServer(useCases usecase.Cases, router *gin.Engine) *Server {
	r := gin.New()
	r.Use(gin.Recovery())

	s := &Server{
		Router:       router,
		MetricRouter: r,
		HttpServer: http.Server{
			Addr:    fmt.Sprintf(":%d", config.C.Server.Port),
			Handler: router,
		},
		MetricHttpServer: http.Server{
			Addr:    fmt.Sprintf(":%d", config.C.Server.MetricsPort),
			Handler: r,
		},
	}

	setupRouter(s, useCases)

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
