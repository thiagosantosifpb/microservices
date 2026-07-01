package ports

import "github.com/thiagosantosifpb/microservices/order/internal/application/core/domain"

type DBPort interface {
	Get(id string) (domain.Order, error)
	Save(order *domain.Order) error
	UpdateStatus(orderID int64, status string) error
	UpdateDeliveryDays(orderID int64, deliveryDays int32) error
	ProductExists(productCode string) (bool, error)
}
