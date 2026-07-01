package ports

import (
	"context"
	"github.com/thiagosantosifpb/microservices/shipping/internal/application/core/domain"
)

type APIPort interface {
	Schedule(ctx context.Context, shipment domain.Shipment) (domain.Shipment, error)
}
