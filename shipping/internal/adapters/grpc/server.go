package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	shippingpb "github.com/thiagosantosifpb/microservices-proto/golang/shipping"
	"github.com/thiagosantosifpb/microservices/shipping/config"
	"github.com/thiagosantosifpb/microservices/shipping/internal/application/core/domain"
	"github.com/thiagosantosifpb/microservices/shipping/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	api  ports.APIPort
	port int
	shippingpb.UnimplementedShippingServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter { return &Adapter{api: api, port: port} }

func (a Adapter) Create(ctx context.Context, request *shippingpb.CreateShippingRequest) (*shippingpb.CreateShippingResponse, error) {
	items := make([]domain.ShippingItem, 0, len(request.GetShippingItems()))
	for _, item := range request.GetShippingItems() {
		items = append(items, domain.ShippingItem{ProductCode: item.GetProductCode(), Quantity: item.GetQuantity()})
	}
	shipment := domain.NewShipment(request.GetOrderId(), items)
	result, err := a.api.Schedule(ctx, shipment)
	code := status.Code(err)
	if code == codes.InvalidArgument {
		return nil, err
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to schedule shipping: %v", err)
	}
	return &shippingpb.CreateShippingResponse{DeliveryDays: result.DeliveryDays}, nil
}

func (a Adapter) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d, error: %v", a.port, err)
	}
	server := grpc.NewServer()
	shippingpb.RegisterShippingServer(server, a)
	if config.GetEnv() == "development" {
		reflection.Register(server)
	}
	log.Printf("Shipping gRPC server running on port %d", a.port)
	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port %d, error: %v", a.port, err)
	}
}
