package place

import (
	"net/http"
	"strconv"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// DeletePlace godoc
// @Summary Delete a place
// @Description Delete a place with same id in the database
// @Tags places
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param id path string true "Place ID"
// @Success 200
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /places/{id} [delete]
func DeletePlace(placeUseCase usecase.Place) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = placeUseCase.DeletePlace(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}
