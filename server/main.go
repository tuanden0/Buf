package main

import (
	"context"
	"fmt"
	"log"
	"net"

	// This import path is based on the name declaration in the go.mod,
	// and the gen/proto/go output location in the buf.gen.yaml.
	userv1 "buf/gen/go/user/v1"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	validate := validator.New()

	// Custom error
	en := en.New()
	uni := ut.New(en, en)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}
	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Fatal(err)
	}

	// Register custom error message
	// required
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.StructField())
		return t
	})

	// email
	validate.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())

		return t
	})

	// role
	validate.RegisterTranslation("role", trans, func(ut ut.Translator) error {
		return ut.Add("role", "{0} only accept 'user' value!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("role", fe.Field())

		return t
	})

	// Register custom valdiate function
	validate.RegisterValidation("role", validateRole)

	// Using validator v10 to validate tags
	err := validate.Struct(in)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		// Return first error, not all error
		return nil, status.Error(codes.InvalidArgument, errs[0].Translate(trans))
	}

	// This validate generate by protoc-gen-validate
	// At this point, only validate email
	// but validator v10 already validate email
	// so this code will run when email is valid
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &userv1.CreateResponse{
		Id:       0,
		Username: in.GetUsername(),
		Email:    in.GetEmail(),
		Role:     in.GetRole(),
	}, nil

}

// Function to validate User Role
func validateRole(fl validator.FieldLevel) bool {
	// Only allow role "user"
	return fl.Field().String() == "user"
}
