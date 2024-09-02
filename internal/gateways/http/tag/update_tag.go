package tag

import (
	"net/http"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// UpdateTag godoc
// @Summary Update a tag
// @Description Update an existing tag in the database
// @Tags places
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param tag body domain.Tag true "Tag data"
// @Success 200 {object} domain.Tag "Updated tag"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /places/tag/{id} [put]
func UpdateTag(tagUseCase usecase.Tag) gin.HandlerFunc {
	return func(c *gin.Context) {
		tag := new(domain.Tag)
		err := c.BindJSON(&tag)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tag, err = tagUseCase.UpdateTag(c, tag)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tag)
	}
}
