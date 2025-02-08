package metric

import (
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metric interface {
	Handle(c *gin.Context)
}

func AllMetrics(h http.Handler, metrics []Metric) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, m := range metrics {
			m.Handle(c)
		}
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	metrics := make([]Metric, 0)
	metrics = append(metrics, NewActiveRoomMetric(cases.RoomRepo))
	metrics = append(metrics, NewSwipeMetric(cases.Swipe))
	h := promhttp.Handler()
	r.GET("/metrics", AllMetrics(h, metrics))
}
