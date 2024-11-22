package place

import (
	"errors"
	"net/http"

	"dishdash.ru/internal/repo"
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// GetPlaceByUrl godoc
// @Summary Get place by url
// @Description Get a place from the database by url
// @Tags places
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param url query string true "place url"
// @Success 200 {object} usecase.Place "place data"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /place/by_url [get]
func GetPlaceByUrl(placeUseCase usecase.Place) gin.HandlerFunc {
	return func(c *gin.Context) {
		place, err := placeUseCase.GetPlaceByUrl(c, c.Query("url"))
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
