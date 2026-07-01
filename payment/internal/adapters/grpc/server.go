package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	paymentpb "github.com/thiagosantosifpb/microservices-proto/golang/payment"
	"github.com/thiagosantosifpb/microservices/payment/config"
	"github.com/thiagosantosifpb/microservices/payment/internal/application/core/domain"
	"github.com/thiagosantosifpb/microservices/payment/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	api  ports.APIPort
	port int
	paymentpb.UnimplementedPaymentServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter { return &Adapter{api: api, port: port} }

func (a Adapter) Create(ctx context.Context, request *paymentpb.CreatePaymentRequest) (*paymentpb.CreatePaymentResponse, error) {
	log.Printf("Creating payment for order %d", request.GetOrderId())
	newPayment := domain.NewPayment(request.GetUserId(), request.GetOrderId(), request.GetTotalPrice())
	result, err := a.api.Charge(ctx, newPayment)
	code := status.Code(err)
	if code == codes.InvalidArgument {
		return nil, err
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to charge: %v", err)
	}
	return &paymentpb.CreatePaymentResponse{PaymentId: result.ID, BillId: result.BillID}, nil
}

func (a Adapter) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d, error: %v", a.port, err)
	}
	server := grpc.NewServer()
	paymentpb.RegisterPaymentServer(server, a)
	if config.GetEnv() == "development" {
		reflection.Register(server)
	}
	log.Printf("Payment gRPC server running on port %d", a.port)
	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port %d, error: %v", a.port, err)
	}
}
