package lobby

import (
	"net/http"

	"dishdash.ru/pkg/domain"
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// NearestLobby godoc
// @Summary find nearest lobby
// @Description find nearest lobby in the database
// @Tags lobbies
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param location body domain.Coordinate true "Location"
// @Success 200 {object} nearestLobbyOutput "Nearest lobby + Distance (in metres)"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /lobbies/nearest [post]
func NearestLobby(lobbyUseCase usecase.Lobby) gin.HandlerFunc {
	return func(c *gin.Context) {
		var location domain.Coordinate
		err := c.BindJSON(&location)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		lobby, dist, err := lobbyUseCase.NearestActiveLobby(c, location)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, nearestLobbyOutput{
			Dist:  dist,
			Lobby: lobby,
		})
	}
}
