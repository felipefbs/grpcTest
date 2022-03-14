package services

import (
	"context"

	"github.com/felipefbs/grpc/pb"
)

type userService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *userService {
	return &userService{}
}

func (u *userService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	return &pb.User{
		Id:    "0",
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}, nil
}
