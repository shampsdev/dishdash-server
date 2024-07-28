package place

import (
	"dishdash.ru/internal/domain"
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// CreateTag godoc
// @Summary Create a tag
// @Description Create a new tag in the database
// @Tags places
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param tag body domain.Tag true "Tag data"
// @Success 200 {object} domain.Tag "Saved tag"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /places/tags [post]
func CreateTag(tagUseCase usecase.Tag) gin.HandlerFunc {
	return func(c *gin.Context) {
		tag := new(domain.Tag)
		err := c.BindJSON(&tag)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tag, err = tagUseCase.SaveTag(c, tag)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tag)
	}
}
