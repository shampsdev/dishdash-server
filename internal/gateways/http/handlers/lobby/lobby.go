package lobby

import (
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, lobbyUseCase *usecase.Lobby) {
	cardGroup := r.Group("lobby")
	cardGroup.POST("", SaveCard(lobbyUseCase))
}
