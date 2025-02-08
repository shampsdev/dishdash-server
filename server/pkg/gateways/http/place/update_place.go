package place

import (
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// UpdatePlace godoc
// @Summary Update a place
// @Description Update a place with same id in the database
// @Tags places
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param place body usecase.UpdatePlaceInput true "Place data"
// @Success 200 {object} domain.Place "Updated place"
// @Failure 400 "Bad Request"
// @Failure 500 "pkg Server Error"
// @Security ApiKeyAuth
// @Router /places [put]
func UpdatePlace(placeUseCase usecase.Place) gin.HandlerFunc {
	return func(c *gin.Context) {
		var placeInput usecase.UpdatePlaceInput
		err := c.BindJSON(&placeInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		place, err := placeUseCase.UpdatePlace(c, placeInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, place)
	}
}
