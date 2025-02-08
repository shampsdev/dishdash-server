package tag

import (
	"net/http"
	"strconv"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// DeleteTag godoc
// @Summary Delete a tag
// @Description Delete an existing tag from the database
// @Tags places
// @Param id path int true "Tag ID"
// @Success 200 "Tag deleted"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /places/tag/{id} [delete]
func DeleteTag(tagUseCase usecase.Tag) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
			return
		}

		err = tagUseCase.DeleteTag(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}
