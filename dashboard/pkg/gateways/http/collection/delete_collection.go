package collection

import (
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// DeleteCollection godoc
// @Summary Delete a collection
// @Description Delete a collection with same id in the database
// @Tags collections
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param id path string true "Collection ID"
// @Success 200
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /collections/{id} [delete]
func DeleteCollection(collectionUseCase usecase.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

        err := collectionUseCase.DeleteCollection(c, id)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        c.Status(http.StatusOK)
    }
}
