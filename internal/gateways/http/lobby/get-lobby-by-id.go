package lobby

import (
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// GetLobbyByID godoc
// @Summary Get lobby by ID
// @Description Get a lobby from the database by ID
// @Tags lobbies
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param id path string true "Lobby ID"
// @Success 200 {object} lobbyOutput "Lobby data"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /lobbies/{id} [get]
func GetLobbyByID(lobbyUseCase usecase.Lobby) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		lobby, err := lobbyUseCase.GetLobbyByID(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, lobbyToOutput(lobby))
	}
}
