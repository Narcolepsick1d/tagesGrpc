package handler

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/Narcolepsick1d/testTages/internal/repository"
	pb "github.com/Narcolepsick1d/testTages/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var ErrNoImage = errors.New("no images")

const (
	maxUploads = 10
	maxViews   = 100
)

type ImageServer struct {
	image   repository.ImageService
	uploads int
	views   int
	sync.Mutex
	pb.UnimplementedImageServiceServer
}

func NewImageServer(image repository.ImageService) *ImageServer {
	return &ImageServer{
		image: image,
	}
}

func (i *ImageServer) DownloadImage(ctx context.Context, req *pb.ImageRequest) (*pb.ImageData, error) {
	i.Lock()
	if i.uploads >= maxUploads {
		return nil, status.Error(codes.Aborted, "too many concurrent uploads")
	}
	i.uploads++
	i.Unlock()
	defer func() {
		i.Lock()
		i.uploads--
		i.Unlock()
	}()
	data, err := i.image.DownloadImage(ctx, req.ImagesName)
	if err != nil {
		if errors.Is(err, ErrNoImage) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pb.ImageData{Data: data, ImageName: req.ImagesName}, nil
}

func (i *ImageServer) UploadImage(ctx context.Context, req *pb.ImageData) (*emptypb.Empty, error) {
	i.Lock()
	if i.uploads >= maxUploads {
		return nil, status.Error(codes.Aborted, "too many concurrent uploads")
	}
	i.uploads++
	defer func() {
		i.Lock()
		i.uploads--
		i.Unlock()
	}()
	err := i.image.UploadImage(ctx, req.ImageName, req.Data)
	if err != nil {
		log.Print(err)
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (i *ImageServer) ListImages(req *emptypb.Empty, srv pb.ImageService_ListImagesServer) error {
	i.Lock()
	if i.views >= maxViews {
		return status.Error(codes.Aborted, "too many concurrent images")
	}
	i.views++
	i.Unlock()
	defer func() {
		i.Lock()
		i.views--
		i.Unlock()
	}()
	images, err := i.image.ListImages(context.Background())
	if err != nil {
		log.Print(err)
		return status.Error(codes.Unknown, err.Error())
	}

	for _, image := range images {
		srv.Send(&pb.ImageInfo{
			Name:       image.Name,
			CreateDate: image.CreatedTime.String(),
			UpdateDate: image.UpdatedTime.String(),
		})
	}

	return nil
}
