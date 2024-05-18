package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter(r *gin.Engine) {
	r.HandleMethodNotAllowed = true
	r.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})
}
