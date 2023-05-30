package main

import (
	"context"
	"github.com/Narcolepsick1d/testTages/handler"
	"github.com/Narcolepsick1d/testTages/internal/service"
	pb "github.com/Narcolepsick1d/testTages/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

const port = ":5001"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	image := service.NewImage()
	imageService := handler.NewImageServer(image)
	serverGRPC := grpc.NewServer()
	pb.RegisterImageServiceServer(serverGRPC, imageService)
	go func() {
		listen, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Server started at port:%s", port)
		err = serverGRPC.Serve(listen)
		if err != nil {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
}
