package collection

import (
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// SaceCollection godoc
// @Summary Create a collection
// @Description Create a new collection in the database
// @Tags collections
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param collection body usecase.SaveCollectionInput true "Collection data"
// @Success 200 {object} domain.Collection "Saved collection"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /collections [post]
func SaveCollection(collectionUseCase usecase.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
        var collectionInput usecase.SaveCollectionInput
        err := c.BindJSON(&collectionInput)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        collection, err := collectionUseCase.SaveCollection(c, collectionInput)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, collection)
    }
}