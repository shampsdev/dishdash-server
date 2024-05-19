package lobby

import (
	"net/http"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/dto"
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// SaveCard godoc
// @Summary Save a lobby
// @Description Save a new lobby to the database
// @Tags lobby
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param card body dto.LobbyToCreate true "Lobby data"
// @Success 200 {object} dto.Lobby "Saved lobby"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /lobby [post]
func SaveCard(lobbyUseCase *usecase.Lobby) gin.HandlerFunc {
	return func(c *gin.Context) {
		var lobbyDto dto.LobbyToCreate
		err := c.BindJSON(&lobbyDto)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		lobby, err := domain.LobbyFromDtoToCreate(lobbyDto)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err = lobbyUseCase.SaveLobby(c, lobby)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusCreated, domain.LobbyToDto(*lobby))
	}
}
