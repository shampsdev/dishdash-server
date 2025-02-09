package metric

import (
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type SwipeMetric struct {
	metric prometheus.Gauge
	swipe  usecase.Swipe
}

func NewSwipeMetric(useCase usecase.Swipe) *SwipeMetric {
	swipeCountMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "swipes_count",
			Help: "A count of swipes",
		})
	prometheus.MustRegister(swipeCountMetric)

	return &SwipeMetric{
		metric: swipeCountMetric,
		swipe:  useCase,
	}
}

func (s SwipeMetric) Handle(c *gin.Context) {
	swipeCount, err := s.swipe.GetCount(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	s.metric.Set(float64(swipeCount))
}
