package tag

import (
	"net/http"

	"dishdash.ru/pkg/usecase"
	"github.com/gin-gonic/gin"
)

// GetAllTags godoc
// @Summary Get tags
// @Description Get a list of tags from the database
// @Tags places
// @Accept  json
// @Produce  json
// @Schemes http https
// @Success 200 {array} domain.Tag "List of tags"
// @Failure 500
// @Router /places/tag [get]
func GetAllTags(tagUseCase usecase.Tag) gin.HandlerFunc {
	return func(c *gin.Context) {
		tags, err := tagUseCase.GetAllTags(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, tags)
	}
}
