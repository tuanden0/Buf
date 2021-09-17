package main

import (
	"context"
	"fmt"
	"log"

	// This import path is based on the name declaration in the go.mod,
	// and the gen/proto/go output location in the buf.gen.yaml.
	userv1 "buf/gen/go/user/v1"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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
		Username: "",
		Password: "",
		Email:    "test",
		Role:     "admin",
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			for _, detail := range e.Details() {
				switch t := detail.(type) {
				case *errdetails.BadRequest:
					for _, violation := range t.GetFieldViolations() {
						fmt.Printf("The %q field was wrong:\n", violation.GetField())
						fmt.Printf("\t%s\n", violation.GetDescription())
					}
				}
			}
		}

		return fmt.Errorf("server validate error: %v", err)
	}

	log.Println("successfully Create", u)
	return nil
}
