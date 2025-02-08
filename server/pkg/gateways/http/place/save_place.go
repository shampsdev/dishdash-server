package place

import (
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// SavePlace godoc
// @Summary Create a place
// @Description Create a new place in the database
// @Tags places
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param place body usecase.SavePlaceInput true "Place data"
// @Success 200 {object} domain.Place "Saved place"
// @Failure 400 "Bad Request"
// @Failure 500 "pkg Server Error"
// @Security ApiKeyAuth
// @Router /places [post]
func SavePlace(placeUseCase usecase.Place) gin.HandlerFunc {
	return func(c *gin.Context) {
		var placeInput usecase.SavePlaceInput
		err := c.BindJSON(&placeInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		place, err := placeUseCase.SavePlace(c, placeInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, place)
	}
}
