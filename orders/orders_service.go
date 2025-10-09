package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/odyssey121/proto_services_go/protos/golang/orders"
	"google.golang.org/grpc"
)

const (
	responseOk  = "ok"
	responseErr = "error"
	serverAddr  = ":50556"
)

type OrdersServiceServer struct {
	orders.UnimplementedOrderServiceServer	
}

func (s *OrdersServiceServer) PlaceOrder(ctx context.Context, req *orders.PlaceOrderRequest) (*orders.PlaceOrderResponse, error) {
	OrderID := "1"
	log.Printf("Processing place order for user: %s, items: %v, payment method: %s", req.UserId, req.Items, req.PaymentMethod)
	return &orders.PlaceOrderResponse{Success: true, OrderId: OrderID, Message: responseOk}, nil

}

func main() {
	lis, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	orders.RegisterOrderServiceServer(grpcServer, &OrdersServiceServer{})

	fmt.Printf("Start gRPC server on %s\n", serverAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server on %s\n", serverAddr)
	}

}
