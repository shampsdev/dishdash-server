package collection

import (
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// GetCollectionByID godoc
// @Summary Get a collection preview
// @Description Get a collection preview with same id from database
// @Tags collections
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param id path string true "Collection ID"
// @Success 200 {object} domain.CollectionPreview "Collection"
// @Failure 500
// @Security ApiKeyAuth
// @Router /collections/preview/{id} [get]
func GetCollectionPreviewByID(collectionUseCase usecase.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		collection, err := collectionUseCase.GetCollectionPreviewByID(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, collection)
	}
}
