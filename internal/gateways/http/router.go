package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func setupRouter(r *gin.Engine) {
	r.HandleMethodNotAllowed = true
	r.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})
}
