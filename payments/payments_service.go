package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/odyssey121/proto_services_go/protos/golang/payments"
	"google.golang.org/grpc"
)

const (
	responseOk  = "ok"
	responseErr = "error"
	ServerAddr  = ":50555"
)

type PaymentsServiceServer struct {
	payments.UnimplementedPaymentServiceServer
}

func (s *PaymentsServiceServer) MakePayment(ctx context.Context, req *payments.MakePaymentRequest) (*payments.MakePaymentResponse, error) {
	requestID := uuid.New().String()
	log.Printf("Processing payment for user: %s, amount: %f, requestID: %s", req.UserId, req.Amount, requestID)
	return &payments.MakePaymentResponse{Success: true, Uuid: requestID, Message: responseOk}, nil

}

func main() {
	lis, err := net.Listen("tcp", ServerAddr)
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	payments.RegisterPaymentServiceServer(grpcServer, &PaymentsServiceServer{})

	fmt.Printf("Start gRPC server on %s\n", ServerAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server on %s\n", ServerAddr)
	}

}
