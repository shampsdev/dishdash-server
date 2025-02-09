package photo

import (
	"net/http"

	"dashboard.dishdash.ru/cmd/config"
	"dashboard.dishdash.ru/pkg/photo"
	"github.com/gin-gonic/gin"
)

type UploadResponse struct {
	URL string `json:"url"`
}

// Upload godoc
// @Summary Upload a image by url to s3
// @Description Upload a image by url to s3
// @Tags photo
// @Accept json
// @Produce json
// @Schemes http https
// @Param url query string true "Url"
// @Success 200 {object} UploadResponse "A url to the stored image"
// @Failure 400 "Parsing error"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /photo/upload [post]
func Upload(cfg config.S3Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		imageURL := c.Query("url")
		if imageURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No URL provided"})
			return
		}

		url, err := photo.UploadToS3(cfg, imageURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"url": url})
	}
}
