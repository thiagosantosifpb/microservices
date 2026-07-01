package api

import (
	"fmt"

	"github.com/thiagosantosifpb/microservices/order/internal/application/core/domain"
	"github.com/thiagosantosifpb/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db       ports.DBPort
	payment  ports.PaymentPort
	shipping ports.ShippingPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort, shipping ports.ShippingPort) *Application {
	return &Application{db: db, payment: payment, shipping: shipping}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	if len(order.OrderItems) == 0 {
		return domain.Order{}, status.Error(codes.InvalidArgument, "order must contain at least one item")
	}
	if order.TotalQuantity() > 50 {
		return domain.Order{}, status.Errorf(codes.InvalidArgument, "orders with more than 50 items are not allowed. total quantity: %d", order.TotalQuantity())
	}
	for _, item := range order.OrderItems {
		if item.Quantity <= 0 {
			return domain.Order{}, status.Errorf(codes.InvalidArgument, "item %s has invalid quantity", item.ProductCode)
		}
		exists, err := a.db.ProductExists(item.ProductCode)
		if err != nil {
			return domain.Order{}, status.Errorf(codes.Internal, "failed to validate product %s: %v", item.ProductCode, err)
		}
		if !exists {
			return domain.Order{}, status.Errorf(codes.NotFound, "product %s not found in stock", item.ProductCode)
		}
	}

	if err := a.db.Save(&order); err != nil {
		return domain.Order{}, status.Errorf(codes.Internal, "failed to save order: %v", err)
	}
	if err := a.payment.Charge(&order); err != nil {
		_ = a.db.UpdateStatus(order.ID, "Canceled")
		return domain.Order{}, err
	}
	if err := a.db.UpdateStatus(order.ID, "Paid"); err != nil {
		return domain.Order{}, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}
	order.Status = "Paid"

	deliveryDays, err := a.shipping.Schedule(&order)
	if err != nil {
		_ = a.db.UpdateStatus(order.ID, "Canceled")
		return domain.Order{}, fmt.Errorf("failed to schedule shipping: %w", err)
	}
	order.DeliveryDays = deliveryDays
	if err := a.db.UpdateDeliveryDays(order.ID, deliveryDays); err != nil {
		return domain.Order{}, status.Errorf(codes.Internal, "failed to update delivery days: %v", err)
	}
	return order, nil
}
