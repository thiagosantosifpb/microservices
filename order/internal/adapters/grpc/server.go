package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	orderpb "github.com/thiagosantosifpb/microservices-proto/golang/order"
	"github.com/thiagosantosifpb/microservices/order/config"
	"github.com/thiagosantosifpb/microservices/order/internal/application/core/domain"
	"github.com/thiagosantosifpb/microservices/order/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	api  ports.APIPort
	port int
	orderpb.UnimplementedOrderServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter { return &Adapter{api: api, port: port} }

func (a Adapter) Create(ctx context.Context, request *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	items := make([]domain.OrderItem, 0, len(request.GetOrderItems()))
	for _, item := range request.GetOrderItems() {
		items = append(items, domain.OrderItem{ProductCode: item.GetProductCode(), UnitPrice: item.GetUnitPrice(), Quantity: item.GetQuantity()})
	}
	newOrder := domain.NewOrder(int64(request.GetCostumerId()), items)
	result, err := a.api.PlaceOrder(newOrder)
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Errorf(codes.Internal, "failed to place order: %v", err)
	}
	return &orderpb.CreateOrderResponse{OrderId: int32(result.ID)}, nil
}

func (a Adapter) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d, error: %v", a.port, err)
	}
	grpcServer := grpc.NewServer()
	orderpb.RegisterOrderServer(grpcServer, a)
	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}
	log.Printf("Order gRPC server running on port %d", a.port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port %d, error: %v", a.port, err)
	}
}
