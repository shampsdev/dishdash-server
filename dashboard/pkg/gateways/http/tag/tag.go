package tag

import (
	"dashboard.dishdash.ru/cmd/config"
	"dashboard.dishdash.ru/pkg/gateways/http/middlewares"
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	tagGroup := r.Group("places").Group("tag")
	tagGroup.Use(middlewares.ApiTokenAuth(config.C.Auth.ApiToken))

	tagGroup.
		GET("", GetAllTags(cases.Tag)).
		POST("", CreateTag(cases.Tag)).
		PUT("", UpdateTag(cases.Tag)).
		DELETE("id/:id", DeleteTag(cases.Tag))
}
