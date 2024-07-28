package lobby

import (
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// SaveLobby godoc
// @Summary Create a lobby
// @Description Create a new lobby in the database
// @Tags lobbies
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param lobby body usecase.SaveLobbyInput true "Lobby data"
// @Success 200 {object} domain.Lobby "Saved lobby"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
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
