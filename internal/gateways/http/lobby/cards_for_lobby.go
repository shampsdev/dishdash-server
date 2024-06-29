package lobby

import (
	"net/http"

	"dishdash.ru/pkg/filter"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// CardsForLobby godoc
// @Summary Get cards filtered with lobby settings
// @Description Get cards filtered with lobby settings
// @Tags lobbies
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param id path string true "Lobby ID"
// @Success 200 {array} cardOutput "cards"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /lobbies/{id}/cards [get]
func CardsForLobby(lobbyUseCase usecase.Lobby) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		lobby, err := lobbyUseCase.GetLobbyByID(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cards, err := lobbyUseCase.GetCardsForSettings(c, lobby.Location, lobby.LobbySettings)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, filter.Map(cards, cardToOutput))
	}
}
