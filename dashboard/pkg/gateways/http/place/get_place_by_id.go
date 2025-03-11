package place

import (
	"errors"
	"net/http"
	"strconv"

	"dishdash.ru/pkg/repo"
	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// GetPlaceByID godoc
// @Summary Get place by id
// @Description Get a place from the database by id
// @Tags places
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param id path string true "Place ID"
// @Success 200 {object} usecase.Place "place data"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /places/id/{id} [get]
func GetPlaceByID(placeUseCase usecase.Place) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		place, err := placeUseCase.GetPlaceByID(c, id)
		if errors.Is(err, repo.ErrPlaceNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, place)
	}
}
