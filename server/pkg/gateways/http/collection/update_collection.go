package collection

import (
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// UpdateCollection godoc
// @Summary Update a collection
// @Description Update a collection with same id in the database
// @Tags collections
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param collection body usecase.UpdateCollectionInput true "Collection data"
// @Success 200 {object} domain.Collection "Updated collection"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /collections [put]
func UpdateCollection(collectionUseCase usecase.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var collectionInput usecase.UpdateCollectionInput
		err := c.BindJSON(&collectionInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		collection, err := collectionUseCase.UpdateCollection(c, collectionInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, collection)
	}
}
