package lobby

import (
	"net/http"

	"dishdash.ru/pkg/domain"
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// CreateLobby godoc
// @Summary Create a lobby with given settings
// @Description Create a new lobby in the database with given settings
// @Tags lobbies
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param lobby body domain.LobbySettings true "lobby settings"
// @Success 200 {object} domain.Lobby "Saved lobby"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /lobbies [post]
func SaveLobby(lobbyUseCase usecase.Lobby) gin.HandlerFunc {
	return func(c *gin.Context) {
		var settings domain.LobbySettings
		err := c.BindJSON(&settings)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		lobby, err := lobbyUseCase.CreateLobby(c, settings)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, lobby)
	}
}
