package lobby

import (
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	lobbiesGroup := r.Group("lobbies")
	lobbiesGroup.POST("", CreateLobby(cases.Lobby))
	lobbiesGroup.GET("/:id", GetLobbyByID(cases.Lobby))
	lobbiesGroup.DELETE("/:id", DeleteLobby(cases.Lobby))
	lobbiesGroup.POST("/nearest", NearestLobby(cases.Lobby))
	lobbiesGroup.POST("/find", FindLobby(cases.Lobby))
}
