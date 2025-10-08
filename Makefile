proto-orders:
	mkdir -p protos/golang/orders
	rm -rfv protos/golang/orders/*
	protoc --proto_path=protos --go_out=protos/golang/orders --go_opt=paths=source_relative \
    --go-grpc_out=protos/golang/orders --go-grpc_opt=paths=source_relative \
    protos/order.proto
	cd protos/golang/orders && go mod init github.com/odyssey121/proto_services_go/protos/golang/orders && go mod tidy
	
proto-payments:
	mkdir -p protos/golang/payments
	rm -rfv protos/golang/payments/*
	protoc --proto_path=protos --go_out=protos/golang/payments --go_opt=paths=source_relative \
    --go-grpc_out=protos/golang/payments --go-grpc_opt=paths=source_relative \
    protos/payment.proto
	cd protos/golang/payments && go mod init github.com/odyssey121/proto_services_go/protos/golang/payments && go mod tidy

# not file in dir
.PHONY: proto-orders proto-payments