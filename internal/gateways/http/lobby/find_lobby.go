package lobby

import (
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// FindLobby godoc
// @Summary find lobby
// @Description shortcut for find nearest + create if not close enough
// @Tags lobbies
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param location body usecase.FindLobbyInput true "Location + Distance (in metres)"
// @Success 200 {object} lobbyOutput
// @Success 201 {object} lobbyOutput
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /lobbies/find [post]
func FindLobby(lobbyUseCase usecase.Lobby) gin.HandlerFunc {
	return func(c *gin.Context) {
		var locDist usecase.FindLobbyInput
		err := c.BindJSON(&locDist)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		lobby, err := lobbyUseCase.FindLobby(c.Request.Context(), locDist)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		lobbyOutput := ToLobbyOutput(lobby)

		c.JSON(http.StatusOK, lobbyOutput)
	}
}
