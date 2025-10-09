package main

import (
	"context"
	"fmt"

	"github.com/odyssey121/proto_services_go/protos/golang/discounts"
)

const (
	responseOk  = "ok"
	responseErr = "error"
	ServerAddr  = ":50557"
)

type DiscountServiceServer struct {
	discounts.UnimplementedDiscountServiceServer	
}

func (s *DiscountServiceServer) MakeDiscount(ctx context.Context, req discounts.MakeDiscountRequest) (*discounts.MakeDiscountResponse, error) {

}

func main() {
	fmt.Println("discounts", discounts.DiscountItem{})

}
