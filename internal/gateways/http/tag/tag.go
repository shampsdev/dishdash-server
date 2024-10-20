package tag

import (
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	placeGroup := r.Group("places").Group("tags")
	placeGroup.POST("", CreateTag(cases.Tag))
	placeGroup.GET("", GetAllTags(cases.Tag))
	placeGroup.PUT("", UpdateTag(cases.Tag))
	placeGroup.DELETE(":id", DeleteTag(cases.Tag))
}
