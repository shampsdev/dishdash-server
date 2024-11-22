package metric

import (
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type ActiveRoomMetric struct {
	metric   prometheus.Gauge
	roomRepo usecase.RoomRepo
}

func NewActiveRoomMetric(repo usecase.RoomRepo) *ActiveRoomMetric {
	activeRoomCount := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_room_count",
			Help: "Total number of active rooms.",
		},
	)
	prometheus.MustRegister(activeRoomCount)

	return &ActiveRoomMetric{
		metric:   activeRoomCount,
		roomRepo: repo,
	}
}

func (a *ActiveRoomMetric) Handle(c *gin.Context) {
	count, err := a.roomRepo.GetActiveRoomCount()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a.metric.Set(float64(count))
}
