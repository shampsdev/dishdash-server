package tag

import (
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	placeGroup := r.Group("places").Group("tag")

	placeGroup.GET("", GetAllTags(cases.Tag))
}
