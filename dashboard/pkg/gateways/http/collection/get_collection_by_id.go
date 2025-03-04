package collection

import (
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// GetCollectionByID godoc
// @Summary Get a collection
// @Description Get a collection with same id from database
// @Tags collections
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param id path string true "Collection ID"
// @Success 200 {object} domain.Collection "Collection"
// @Failure 500
// @Security ApiKeyAuth
// @Router /collections/{id} [get]
func GetCollectionByID(collectionUseCase usecase.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		collection, err := collectionUseCase.GetCollectionByID(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, collection)
	}
}
