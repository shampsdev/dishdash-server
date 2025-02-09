package place

import (
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetAllPlaces godoc
// @Summary Get places
// @Description Get a list of places from the database
// @Tags places
// @Accept  json
// @Produce  json
// @Schemes http https
// @Success 200 {array} domain.Place "List of places"
// @Failure 500
// @Security ApiKeyAuth
// @Router /places [get]
func GetAllPlaces(placeUseCase usecase.Place) gin.HandlerFunc {
	return func(c *gin.Context) {
		places, err := placeUseCase.GetAllPlaces(c)
		if err != nil {
			log.WithError(err).Error("failed to get places")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, places)
	}
}
