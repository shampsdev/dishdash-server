package lobby

import (
	"errors"
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
// @Param location body findLobbyInput true "Location + Distance (in metres)"
// @Success 200 {object} lobbyOutput
// @Success 201 {object} lobbyOutput
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /lobbies/find [post]
func FindLobby(lobbyUseCase usecase.Lobby) gin.HandlerFunc {
	return func(c *gin.Context) {
		var locDist findLobbyInput
		err := c.BindJSON(&locDist)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		lobby, dist, err := lobbyUseCase.NearestActiveLobby(c, locDist.Location)
		if err != nil && !errors.Is(err, usecase.ErrLobbyNotFound) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if dist <= locDist.Dist && !errors.Is(err, usecase.ErrLobbyNotFound) {
			c.JSON(http.StatusOK, lobbyToOutput(lobby))
			return
		}

		lobby, err = lobbyUseCase.CreateLobby(c, usecase.SaveLobbyInput{Location: locDist.Location})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, lobbyToOutput(lobby))
	}
}
