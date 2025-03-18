package place

import (
	"net/http"

	"dishdash.ru/pkg/repo"
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// FilterPlaces godoc
// @Summary Filter places
// @Tags places
// @Accept json
// @Produce json
// @Schemes http https
// @Param place body repo.PlacesFilter true "Filter params"
// @Success 200 {object} []domain.Place "Matched places"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /places/filter [post]
func FilterPlaces(placeUseCase usecase.Place) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter repo.PlacesFilter
		if err := c.ShouldBindJSON(&filter); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		places, err := placeUseCase.FilterPlaces(c, filter)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, places)
	}
}
