package image

import (
	"net/http"

	"dashboard.dishdash.ru/pkg/repo"
	"github.com/gin-gonic/gin"
)

type UploadResponse struct {
	URL string `json:"url"`
}

type uploadByURLRequest struct {
	URL       string `json:"url"`
	Directory string `json:"directory"`
}

// UploadByURL godoc
// @Summary UploadByURL a image by url to s3
// @Description UploadByURL a image by url to s3
// @Tags images
// @Accept json
// @Produce json
// @Schemes http https
// @Param uploadByURLRequest body uploadByURLRequest true "URL and directory in s3 storage"
// @Success 200 {object} UploadResponse "A url to the stored image"
// @Failure 400 "Parsing error"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Security ApiKeyAuth
// @Router /images/upload/by_url [post]
func UploadByURL(storage repo.ImageStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req uploadByURLRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		url, err := storage.SaveImageByURL(c, req.URL, req.Directory)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"url": url})
	}
}
