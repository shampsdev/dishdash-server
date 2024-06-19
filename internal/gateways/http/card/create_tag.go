package card

import (
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// CreateTag godoc
// @Summary Create a tag
// @Description Create a new tag in the database
// @Tags cards
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param tag body tagInput true "Tag data"
// @Success 200 {object} tagOutput "Saved tag"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /cards/tags [post]
func CreateTag(tagUseCase usecase.Tag) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tagInput usecase.TagInput
		err := c.BindJSON(&tagInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		tag, err := tagUseCase.CreateTag(c, tagInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, tagToOutput(tag))
	}
}
