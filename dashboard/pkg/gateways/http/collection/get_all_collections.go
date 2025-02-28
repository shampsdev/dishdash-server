package collection

import (
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetAllCollections godoc
// @Summary Get collections
// @Description Get a list of collections from the database
// @Tags collections
// @Accept  json
// @Produce  json
// @Schemes http https
// @Success 200 {array} domain.Collection "List of collections"
// @Failure 500
// @Security ApiKeyAuth
// @Router /collections [get]
func GetAllCollections(collectionUseCase usecase.Collection) gin.HandlerFunc {
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