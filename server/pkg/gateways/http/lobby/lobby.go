package lobby

import (
	"dishdash.ru/pkg/gateways/http/middlewares"
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	lobbiesGroup := r.Group("lobbies")
	lobbiesGroup.POST("", SaveLobby(cases.Lobby))
	lobbiesGroup.GET("/:id", GetLobbyByID(cases.Lobby))
	lobbiesGroup.POST("/nearest", NearestLobby(cases.Lobby))
	lobbiesGroup.POST("/find", FindLobby(cases.Lobby))

	lobbiesGroupProtected := lobbiesGroup.Group("")
	lobbiesGroupProtected.Use(middlewares.ApiTokenAuth())

	lobbiesGroupProtected.DELETE("/:id", DeleteLobby(cases.Lobby))
}
