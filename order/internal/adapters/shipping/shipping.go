package shipping_adapter

import (
	"context"
	"log"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	shippingpb "github.com/thiagosantosifpb/microservices-proto/golang/shipping"
	"github.com/thiagosantosifpb/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Adapter struct{ shipping shippingpb.ShippingClient }

func NewAdapter(shippingServiceURL string) (*Adapter, error) {
	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
			grpc_retry.WithMax(5),
			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)),
		)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(shippingServiceURL, opts...)
	if err != nil {
		return nil, err
	}
	return &Adapter{shipping: shippingpb.NewShippingClient(conn)}, nil
}

func (a *Adapter) Schedule(order *domain.Order) (int32, error) {
	items := make([]*shippingpb.ShippingItem, 0, len(order.OrderItems))
	for _, item := range order.OrderItems {
		items = append(items, &shippingpb.ShippingItem{ProductCode: item.ProductCode, Quantity: item.Quantity})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	response, err := a.shipping.Create(ctx, &shippingpb.CreateShippingRequest{OrderId: order.ID, ShippingItems: items})
	if status.Code(err) == codes.DeadlineExceeded {
		log.Printf("deadline exceeded while calling Shipping service for order %d", order.ID)
	}
	if err != nil {
		return 0, err
	}
	return response.GetDeliveryDays(), nil
}
