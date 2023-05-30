package service

import (
	"context"
	"errors"
	"github.com/Narcolepsick1d/testTages/internal/model"
	"github.com/Narcolepsick1d/testTages/internal/repository"
	"log"
	"os"
	"path"
	"syscall"
	"time"
)

var ErrNoImage = errors.New("no images")

const Folder = "img/"

type ImageRepository struct {
	imgRepo repository.ImageService
}

func NewImage() *ImageRepository {
	if err := os.Mkdir(Folder, 0755); err != nil {
		log.Print(err)
	}
	log.Println("создана папка img")
	return &ImageRepository{}
}

func (i *ImageRepository) UploadImage(ctx context.Context, imageName string, data []byte) error {
	pathForImg := path.Join(Folder, "/", imageName)
	return os.WriteFile(pathForImg, data, 0755)
}

func (i *ImageRepository) ListImages(ctx context.Context) ([]model.Image, error) {
	log.Println("Started send list files")
	files, err := os.ReadDir(Folder)
	if err != nil {
		log.Fatal(err)
	}
	var imageList []model.Image
	for _, file := range files {
		fileInfo, err := os.Stat(Folder + "/" + file.Name())

		if err != nil {
			log.Fatal(err)
		}
		mTime := fileInfo.ModTime()
		stat := fileInfo.Sys().(*syscall.Win32FileAttributeData)
		creationTime := time.Unix(0, stat.CreationTime.Nanoseconds())
		if !file.IsDir() {
			imageInfo := model.Image{
				Name:        file.Name(),
				CreatedTime: creationTime,
				UpdatedTime: mTime,
			}
			imageList = append(imageList, imageInfo)
		}
	}
	return imageList, nil
}

func (i *ImageRepository) DownloadImage(ctx context.Context, imageName string) ([]byte, error) {
	pathForImg := path.Join(Folder, "/", imageName)
	data, err := os.ReadFile(pathForImg)
	if err != nil {
		if errors.Is(os.ErrNotExist, ErrNoImage) {
			return nil, ErrNoImage
		}

		return nil, err
	}

	return data, nil
}
