package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/odyssey121/proto_services_go/protos/golang/discounts"
	"github.com/odyssey121/proto_services_go/protos/golang/orders"
	"github.com/odyssey121/proto_services_go/protos/golang/payments"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	responseOk  = "ok"
	responseErr = "error"
	ServerAddr  = ":50556"
)

type OrdersServiceServer struct {
	orders.UnimplementedOrderServiceServer
}

func (s *OrdersServiceServer) PlaceOrder(ctx context.Context, req *orders.PlaceOrderRequest) (*orders.PlaceOrderResponse, error) {
	log.Printf("Processing place order for user: %s, items: %v, payment method: %s", req.UserId, req.Items, req.PaymentMethod)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// dicount logic clinet
	discConn, err := grpc.NewClient(":50557", opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to dicount service: %v", err)
	}
	defer discConn.Close()
	// calc total coast
	var totalCoast uint64
	for i := 0; i < len(req.Items); i++ {
		item := req.Items[i]
		totalCoast += uint64(item.GetQuantity()) * uint64((item.GetCoast() * 1000))
	}

	discClinet := discounts.NewDiscountServiceClient(discConn)
	discCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// dicount logic service
	var discResp *discounts.MakeDiscountResponse
	discReq := &discounts.MakeDiscountRequest{UserId: req.UserId, Items: req.Items}
	discResp, err = discClinet.MakeDiscount(discCtx, discReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make discount: %w", err)
	}
	//payment logic client
	paymentConn, err := grpc.NewClient(":50555", opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %v", err)
	}
	defer paymentConn.Close()

	paymentClient := payments.NewPaymentServiceClient(paymentConn)
	paymentCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// payment logic service
	paymentReq := &payments.MakePaymentRequest{UserId: req.UserId, Amount: totalCoast - uint64(discResp.DiscountAmount)}
	paymentRes, err := paymentClient.MakePayment(paymentCtx, paymentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make payment: %w", err)
	}
	// output response logic
	if !paymentRes.Success {
		return &orders.PlaceOrderResponse{
			Success: false,
			OrderId: "",
			Message: fmt.Sprintf("Payment failed, with message: %s, order not created", paymentRes.Message)}, nil
	}

	return &orders.PlaceOrderResponse{Success: true, OrderId: paymentRes.Id, Message: responseOk}, nil

}

func main() {
	lis, err := net.Listen("tcp", ServerAddr)
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	orders.RegisterOrderServiceServer(grpcServer, &OrdersServiceServer{})

	fmt.Printf("Start gRPC server on %s\n", ServerAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server on %s\n", ServerAddr)
	}

}
