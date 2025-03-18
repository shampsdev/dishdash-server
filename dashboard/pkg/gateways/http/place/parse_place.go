package place

import (
	"net/http"

	"dishdash.ru/pkg/domain"
	"github.com/gin-gonic/gin"
)

type ParsePlaceRequest struct {
	Url string `json:"url" binding:"required"`
}

// ParsePlace godoc
// @Summary Parse place with url
// @Tags places
// @Accept json
// @Produce json
// @Schemes http https
// @Param ParsePlaceRequest body ParsePlaceRequest true "Place URL"
// @Success 200 {object} usecase.Place "place data"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /places/parse [post]
func ParsePlace() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ParsePlaceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p := &domain.Place{
			Title:  "Title",
			Url:    &req.Url,
			Images: []string{},
			Tags:   []*domain.Tag{},
		}
		c.JSON(http.StatusOK, p)
	}
}
