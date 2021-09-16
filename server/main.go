package main

import (
	"context"
	"fmt"
	"log"
	"net"

	// This import path is based on the name declaration in the go.mod,
	// and the gen/proto/go output location in the buf.gen.yaml.
	userv1 "buf/gen/go/user/v1"

	"google.golang.org/grpc"
)

// userServiceServer implements the UserService API.
type userServiceServer struct {
	userv1.UnimplementedUserServiceServer
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// Simple and Fast implement GRPC Server
func run() error {
	listenOn := "127.0.0.1:8080"
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", listenOn, err)
	}

	server := grpc.NewServer()
	userv1.RegisterUserServiceServer(server, &userServiceServer{})
	log.Println("Listening on", listenOn)

	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

// Assume implement UserService Create API
func (u *userServiceServer) Create(ctx context.Context, in *userv1.CreateRequest) (*userv1.CreateResponse, error) {

	return &userv1.CreateResponse{
		Id:       0,
		Username: in.GetUsername(),
		Email:    in.GetEmail(),
		Role:     in.GetRole(),
	}, nil

}
