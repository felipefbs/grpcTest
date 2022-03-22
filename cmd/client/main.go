package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/felipefbs/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	AddUsersBiStream(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "123",
		Name:  "felipefbs",
		Email: "a@a.com",
	}
	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "123",
		Name:  "felipefbs",
		Email: "a@a.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(stream)
		fmt.Println("-------------")
		fmt.Printf("Status: %v\n", stream.Status)

	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{Id: "123", Name: "Felipe", Email: "felipe@email.com"},
		{Id: "124", Name: "Henrique", Email: "Henrique@email.com"},
		{Id: "125", Name: "Sara", Email: "Sara@email.com"},
		{Id: "126", Name: "Nicolau", Email: "Nicolau@email.com"},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, req := range reqs {
		fmt.Println(req)
		stream.Send(req)
		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}

func AddUsersBiStream(client pb.UserServiceClient) {
	stream, err := client.AddUserBiStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	reqs := []*pb.User{
		{Id: "123", Name: "Felipe", Email: "felipe@email.com"},
		{Id: "124", Name: "Henrique", Email: "Henrique@email.com"},
		{Id: "125", Name: "Sara", Email: "Sara@email.com"},
		{Id: "126", Name: "Nicolau", Email: "Nicolau@email.com"},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println(req.Name)
			stream.Send(req)
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
				break
			}
			fmt.Println("Recebendo user", res.User.Name, res.Status)
		}
		close(wait)
	}()

	<-wait
}
