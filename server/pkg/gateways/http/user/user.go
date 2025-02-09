package user

import (
	"dishdash.ru/pkg/gateways/http/middlewares"
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

	userGroupProtected := userGroup.Group("")

	userGroupProtected.Use(middlewares.ApiTokenAuth())
	userGroupProtected.GET("", GetAllUsers(cases.User))
}
