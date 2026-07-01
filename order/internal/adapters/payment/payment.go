package payment_adapter

import (
	"context"
	"log"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	paymentpb "github.com/thiagosantosifpb/microservices-proto/golang/payment"
	"github.com/thiagosantosifpb/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Adapter struct{ payment paymentpb.PaymentClient }

func NewAdapter(paymentServiceURL string) (*Adapter, error) {
	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
			grpc_retry.WithMax(5),
			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)),
		)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(paymentServiceURL, opts...)
	if err != nil {
		return nil, err
	}
	return &Adapter{payment: paymentpb.NewPaymentClient(conn)}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := a.payment.Create(ctx, &paymentpb.CreatePaymentRequest{UserId: order.CustomerID, OrderId: order.ID, TotalPrice: order.TotalPrice()})
	if status.Code(err) == codes.DeadlineExceeded {
		log.Printf("deadline exceeded while calling Payment service for order %d", order.ID)
	}
	return err
}
