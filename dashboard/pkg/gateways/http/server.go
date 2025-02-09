package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"dashboard.dishdash.ru/cmd/config"
	"dishdash.ru/pkg/usecase"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/tj/go-spin"

	"golang.org/x/sync/errgroup"
)

const shutdownDuration = 1500 * time.Millisecond

type Server struct {
	HttpServer http.Server
	Router     *gin.Engine
}

func NewServer(useCases usecase.Cases) *Server {
	r := gin.New()
	r.Use(gin.Recovery())

	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.Use(r)

	s := &Server{
		Router: r,
		HttpServer: http.Server{
			Addr:    fmt.Sprintf(":%d", config.C.Server.Port),
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
	err := errors.Join(
		s.HttpServer.Shutdown(ctx),
		eg.Wait(),
	)
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
