package repo

import (
	"context"
	"io"
)

type ImageStorage interface {
	SaveImageByURL(ctx context.Context, url string, destDir string) (string, error)
	SaveImageByBytes(ctx context.Context, imageData []byte, destDir string) (string, error)
	SaveImageByReader(ctx context.Context, imageData io.Reader, destDir string) (string, error)
}
