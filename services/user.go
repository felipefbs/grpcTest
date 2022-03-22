package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

func (*userService) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {
	stream.Send(&pb.UserResultStream{
		Status: "Init",
		User:   &pb.User{},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "Inserting",
		User:   &pb.User{},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "User has been inserted",
		User: &pb.User{
			Id:    "123",
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "Completed",
		User: &pb.User{
			Id:    "123",
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	})

	time.Sleep(time.Second * 3)

	return nil
}

func (u *userService) AddUsers(stream pb.UserService_AddUsersServer) error {
	users := []*pb.User{}

	for {

		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				User: users,
			})
		}
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, &pb.User{
			Id:    req.GetId(),
			Name:  req.GetName(),
			Email: req.GetEmail(),
		})

		fmt.Println("Adding", req.Name)
	}
}
