package lobby

import (
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// DeleteLobby godoc
// @Summary delete a lobby
// @Description delete a lobby in the database
// @Tags lobbies
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param id path string true "lobby id"
// @Success 200
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /lobbies/{id} [delete]
func DeleteLobby(lobbyUseCase usecase.Lobby) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := lobbyUseCase.DeleteLobbyByID(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}
