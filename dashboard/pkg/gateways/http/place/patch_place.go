package place

import (
	"encoding/json"
	"io"
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// PatchPlace godoc
// @Summary Patch a place
// @Tags places
// @Accept json
// @Produce json
// @Schemes http https
// @Param place body usecase.UpdatePlaceInput true "Place data"
// @Success 200 {object} domain.Place "Patched place"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /places [patch]
func PatchPlace(placeUseCase usecase.Place) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var placeInput usecase.UpdatePlaceInput
		err = json.Unmarshal(body, &placeInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		place, err := placeUseCase.GetPlaceByID(c, placeInput.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		placeUpdate := usecase.UpdatePlaceInputFromDomain(place)

		err = json.Unmarshal(body, &placeUpdate)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		placeUpdate.ID = place.ID
		place, err = placeUseCase.UpdatePlace(c, placeUpdate)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, place)
	}
}
