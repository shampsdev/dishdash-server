package photo

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"

	"dashboard.dishdash.ru/cmd/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func UploadToS3(cfg config.S3Config, imageURL string) (string, error) {
	imageData, err := downloadFromURL(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %v", err)
	}

	uniqueImageID := uuid.New()
	mimeType := http.DetectContentType(imageData)
	fileExtension, _ := mime.ExtensionsByType(mimeType)
	if len(fileExtension) == 0 {
		fileExtension = []string{".jpeg"}
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      &cfg.Region,
		Endpoint:    &cfg.EndpointUrl,
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretKey, ""),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	key := fmt.Sprintf("place/%s%s", uniqueImageID, fileExtension[0])
	s3Client := s3.New(sess)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      &cfg.Bucket,
		Key:         aws.String(key),
		Body:        bytes.NewReader(imageData),
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %v", err)
	}

	fileURL := fmt.Sprintf("%s/%s/%s", cfg.EndpointUrl, cfg.Bucket, url.PathEscape(key))
	log.Infof("Image uploaded successfully: %s", fileURL)

	return fileURL, nil
}

func downloadFromURL(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image from URL: %v", err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %v", err)
	}

	return buf.Bytes(), nil
}
