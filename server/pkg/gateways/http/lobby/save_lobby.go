package lobby

import (
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// SaveLobby godoc
// @Summary Create a lobby
// @Description Create a new lobby in the database
// @Tags lobbies
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param lobby body usecase.SaveLobbyInput true "lobby data"
// @Success 200 {object} usecase.LobbyOutput "Saved lobby"
// @Failure 400 "Bad Request"
// @Failure 500 "pkg Server Error"
// @Router /lobbies [post]
func SaveLobby(lobbyUseCase usecase.Lobby) gin.HandlerFunc {
	return func(c *gin.Context) {
		var lobbyInput usecase.SaveLobbyInput
		err := c.BindJSON(&lobbyInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		lobby, err := lobbyUseCase.SaveLobby(c, lobbyInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, lobby)
	}
}
