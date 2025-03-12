package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"

	"dashboard.dishdash.ru/cmd/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

type Storage struct {
	cfg      config.S3Config
	session  *session.Session
	s3Client *s3.S3
}

func NewStorage(cfg config.S3Config) (*Storage, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      &cfg.Region,
		Endpoint:    &cfg.EndpointUrl,
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	s3Client := s3.New(sess)

	return &Storage{
		cfg:      cfg,
		session:  sess,
		s3Client: s3Client,
	}, nil
}

func (s *Storage) SaveImageByURL(_ context.Context, imageURL string, destDir string) (string, error) {
	imageData, err := downloadFromURL(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}

	fileURL, err := s.SaveImageByBytes(context.Background(), imageData, destDir)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	return fileURL, nil
}

func (s *Storage) SaveImageByBytes(_ context.Context, imageData []byte, destDir string) (string, error) {
	uniqueImageID := uuid.New()
	mimeType := http.DetectContentType(imageData)
	fileExtension, _ := mime.ExtensionsByType(mimeType)
	if len(fileExtension) == 0 {
		fileExtension = []string{".jpeg"}
	}

	key := fmt.Sprintf("%s/%s%s", destDir, uniqueImageID, fileExtension[0])

	_, err := s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      &s.cfg.Bucket,
		Key:         aws.String(key),
		Body:        bytes.NewReader(imageData),
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}

	fileURL := fmt.Sprintf("%s/%s/%s", s.cfg.EndpointUrl, s.cfg.Bucket, key)
	log.Infof("Image uploaded successfully: %s", fileURL)

	return fileURL, nil
}

func (s *Storage) SaveImageByReader(_ context.Context, imageData io.Reader, destDir string) (string, error) {
	data, err := io.ReadAll(imageData)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	return s.SaveImageByBytes(context.Background(), data, destDir)
}

func downloadFromURL(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image from URL: %w", err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return buf.Bytes(), nil
}
