package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/odyssey121/proto_services_go/protos/golang/discounts"
	"google.golang.org/grpc"
)

const (
	responseOk  = "ok"
	responseErr = "error"
	ServerAddr  = ":50557"
)

type DiscountServiceServer struct {
	discounts.UnimplementedDiscountServiceServer
}

func (s *DiscountServiceServer) MakeDiscount(ctx context.Context, req *discounts.MakeDiscountRequest) (*discounts.MakeDiscountResponse, error) {
	var totalDiscount uint32
	var totalQuantity int32
	for _, item := range req.Items {
		totalDiscount += item.GetDicount()
		totalQuantity += item.GetQuantity()
	}
	var bonus uint32
	if totalQuantity > 1 {
		bonus = totalDiscount * uint32(totalQuantity) / 100
	}

	totalDiscount += bonus

	return &discounts.MakeDiscountResponse{
		Success:        true,
		DiscountAmount: totalDiscount,
		Message:        responseOk,
	}, nil

}

func main() {
	lis, err := net.Listen("tcp", ServerAddr)
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	discounts.RegisterDiscountServiceServer(grpcServer, &DiscountServiceServer{})

	fmt.Printf("Start gRPC server on %s\n", ServerAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server on %s\n", ServerAddr)
	}

}
