package ports

import "github.com/thiagosantosifpb/microservices/order/internal/application/core/domain"

type PaymentPort interface {
	Charge(order *domain.Order) error
}
