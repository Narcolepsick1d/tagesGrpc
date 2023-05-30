package repository

import (
	"context"
	"github.com/Narcolepsick1d/testTages/internal/model"
)

type ImageService interface {
	UploadImage(ctx context.Context, filename string, data []byte) error
	ListImages(ctx context.Context) ([]model.Image, error)
	DownloadImage(ctx context.Context, filename string) ([]byte, error)
}
