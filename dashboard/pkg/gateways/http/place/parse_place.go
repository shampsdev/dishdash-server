package place

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"dashboard.dishdash.ru/cmd/config"
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

		u, err := url.Parse(fmt.Sprintf("%s/api/parse/", config.C.Parser.URL))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		q := u.Query()
		q.Set("url", req.Url)
		u.RawQuery = q.Encode()

		r := http.Request{
			Method: "GET",
			URL:    u,
			Header: http.Header{
				"api-key": []string{config.C.Parser.ApiKey},
			},
		}

		client := &http.Client{}
		resp, err := client.Do(&r)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		var p domain.Place
		err = json.NewDecoder(resp.Body).Decode(&p)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, p)
	}
}
