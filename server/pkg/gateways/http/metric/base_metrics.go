package metric

import (
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func AddBasicMetrics(r *gin.RouterGroup) {
	p := ginprometheus.NewPrometheus("gin")

	r.Use(p.HandlerFunc())
}
