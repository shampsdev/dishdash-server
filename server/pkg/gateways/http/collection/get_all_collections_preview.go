package collection

import (
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetAllCollectionsPreview godoc
// @Summary Get collections previews
// @Description Get a list of collections preveiws from the database
// @Tags collections
// @Accept  json
// @Produce  json
// @Schemes http https
// @Success 200 {array} domain.CollectionPreview "List of collections previews"
// @Failure 500
// @Router /collections/preview [get]
func GetAllCollectionsPreview(collectionUseCase usecase.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		collections, err := collectionUseCase.GetAllCollections(c)
		if err != nil {
			log.WithError(err).Error("failed to get collections")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, collections)
	}
}
