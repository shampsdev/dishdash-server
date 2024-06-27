package card

import (
	"net/http"

	"dishdash.ru/pkg/filter"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// GetAllTags godoc
// @Summary Get tags
// @Description Get a list of tags from the database
// @Tags cards
// @Accept  json
// @Produce  json
// @Schemes http https
// @Success 200 {array} tagOutput "List of tags"
// @Failure 500
// @Router /cards/tags [get]
func GetAllTags(tagUseCase usecase.Tag) gin.HandlerFunc {
	return func(c *gin.Context) {
		tags, err := tagUseCase.GetAllTags(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, filter.Map(tags, tagToOutput))
	}
}
