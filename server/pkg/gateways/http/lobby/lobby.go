package lobby

import (
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	lobbiesGroup := r.Group("lobbies")
	lobbiesGroup.POST("", SaveLobby(cases.Lobby))
	lobbiesGroup.GET("/:id", GetLobbyByID(cases.Lobby))
}
