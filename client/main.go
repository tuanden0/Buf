package main

import (
	"context"
	"fmt"
	"log"

	// This import path is based on the name declaration in the go.mod,
	// and the gen/proto/go output location in the buf.gen.yaml.
	userv1 "buf/gen/go/user/v1"

	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	connectTo := "127.0.0.1:8080"
	conn, err := grpc.Dial(connectTo, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect to UserService on %s: %w", connectTo, err)
	}
	log.Println("Connected to", connectTo)

	userClient := userv1.NewUserServiceClient(conn)
	u, err := userClient.Create(context.Background(), &userv1.CreateRequest{
		Username: "test",
		Password: "123123",
		Email:    "test",
		Role:     "user",
	})

	if err != nil {

		return fmt.Errorf("server validate error: %v", err)
	}

	log.Println("successfully Create", u)
	return nil
}
