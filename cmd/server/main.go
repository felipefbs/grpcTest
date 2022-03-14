package main

import (
	"log"
	"net"

	"github.com/felipefbs/grpc/pb"
	"github.com/felipefbs/grpc/services"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, services.NewUserService())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
