package lobby

import (
	"dishdash.ru/internal/gateways/http/middlewares"
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	lobbiesGroup := r.Group("lobbies")
	lobbiesGroup.POST("", SaveLobby(cases.Lobby))
	lobbiesGroup.GET("/:id", GetLobbyByID(cases.Lobby))

	lobbiesGroupProtected := lobbiesGroup.Group("")
	lobbiesGroupProtected.Use(middlewares.ApiTokenAuth())

	lobbiesGroupProtected.DELETE("/:id", DeleteLobby(cases.Lobby))
}
