package ports

import (
	"context"
	"github.com/thiagosantosifpb/microservices/payment/internal/application/core/domain"
)

type DBPort interface {
	Save(ctx context.Context, payment *domain.Payment) error
}
