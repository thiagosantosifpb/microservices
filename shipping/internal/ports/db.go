package ports

import (
	"context"
	"github.com/thiagosantosifpb/microservices/shipping/internal/application/core/domain"
)

type DBPort interface {
	Save(ctx context.Context, shipment *domain.Shipment) error
}
