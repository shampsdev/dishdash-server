package user

import (
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	userGroup := r.Group("users")
	userGroup.POST("", SaveUser(cases.User))
	userGroup.POST("with_id", SaveUserWithID(cases.User))
	userGroup.GET(":id", GetUserByID(cases.User))
	userGroup.GET("/telegram/:telegram", GetUserByTelegram(cases.User))
	userGroup.PUT("", UpdateUser(cases.User))
}
