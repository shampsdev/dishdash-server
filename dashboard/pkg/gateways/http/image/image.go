package image

import (
	"dashboard.dishdash.ru/cmd/config"
	"dashboard.dishdash.ru/pkg/gateways/http/middlewares"
	"dashboard.dishdash.ru/pkg/repo/s3"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup) {
	storage, err := s3.NewStorage(config.C.S3)
	if err != nil {
		panic(err)
	}

	imageGroup := r.Group("images")
	imageGroup.Use(middlewares.ApiTokenAuth(config.C.Auth.ApiToken))

	imageGroup.POST("upload/by_url", UploadByURL(storage))
	imageGroup.POST("upload/by_file", UploadByFile(storage))
}
