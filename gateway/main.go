package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	userv1 "buf/gen/go/user/v1"
)

const (
	grpcServer = "127.0.0.1:8080"
	addrStr    = ":8000"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// Handle error
	runtime.HTTPError = CustomHTTPError

	// Connect to GRPC server
	conn, errDial := grpc.DialContext(
		context.Background(), grpcServer,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if errDial != nil {
		return errDial
	}

	// Create MUX
	mux := runtime.NewServeMux()

	// Create UserService Handler
	userHandlerErr := userv1.RegisterUserServiceHandler(context.Background(), mux, conn)
	if userHandlerErr != nil {
		return userHandlerErr
	}

	// Create HTTP Server
	gateway := &http.Server{
		Addr:    addrStr,
		Handler: mux,
	}
	log.Println("HTTP Gateway listening on", addrStr)

	if gatewayErr := gateway.ListenAndServe(); gatewayErr != nil {
		return fmt.Errorf("failed to serve HTTP server: %w", gatewayErr)
	}

	return nil
}

type errorBody struct {
	Code    int32             `json:"code,omitempty"`
	Message string            `json:"message,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// https://mycodesmells.com/post/grpc-gateway-error-handler
func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {

	const fallback = `{"error": "failed to marshal error message"}`

	st := status.Convert(err)
	sc := runtime.HTTPStatusFromCode(status.Code(err))

	errDetails := errorBody{
		Code:    st.Proto().GetCode(),
		Message: grpc.ErrorDesc(err),
		Details: make(map[string]string),
	}

	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *errdetails.BadRequest:
			for _, violation := range t.GetFieldViolations() {
				errDetails.Details[violation.GetField()] = violation.GetDescription()
			}
		}
	}

	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(sc)

	jErr := json.NewEncoder(w).Encode(errDetails)

	if jErr != nil {
		w.Write([]byte(fallback))
	}
}
