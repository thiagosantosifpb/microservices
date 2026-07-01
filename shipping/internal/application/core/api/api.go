package api

import (
	"context"
	"github.com/thiagosantosifpb/microservices/shipping/internal/application/core/domain"
	"github.com/thiagosantosifpb/microservices/shipping/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct{ db ports.DBPort }

func NewApplication(db ports.DBPort) *Application { return &Application{db: db} }

func (a Application) Schedule(ctx context.Context, shipment domain.Shipment) (domain.Shipment, error) {
	if shipment.OrderID <= 0 {
		return domain.Shipment{}, status.Error(codes.InvalidArgument, "order_id is required")
	}
	if len(shipment.Items) == 0 || shipment.TotalQuantity() <= 0 {
		return domain.Shipment{}, status.Error(codes.InvalidArgument, "shipping requires at least one valid item")
	}
	shipment.DeliveryDays = shipment.CalculateDeliveryDays()
	if err := a.db.Save(ctx, &shipment); err != nil {
		return domain.Shipment{}, err
	}
	return shipment, nil
}
