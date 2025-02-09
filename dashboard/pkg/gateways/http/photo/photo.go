package photo

import (
	"dashboard.dishdash.ru/cmd/config"
	"dashboard.dishdash.ru/pkg/gateways/http/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup) {
	photoGroup := r.Group("photo")
	photoGroup.Use(middlewares.ApiTokenAuth(config.C.Auth.ApiToken))

	photoGroup.POST("upload", Upload(config.C.S3))
}
