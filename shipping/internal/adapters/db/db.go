package db

import (
	"context"
	"fmt"

	"github.com/thiagosantosifpb/microservices/shipping/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Shipment struct {
	gorm.Model
	OrderID      int64
	DeliveryDays int32
	Items        []ShippingItem
}

type ShippingItem struct {
	gorm.Model
	ProductCode string
	Quantity    int32
	ShipmentID  uint
}

type Adapter struct{ db *gorm.DB }

func NewAdapter(dataSourceURL string) (*Adapter, error) {
	database, openErr := gorm.Open(mysql.Open(dataSourceURL), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}
	if err := database.AutoMigrate(&Shipment{}, &ShippingItem{}); err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: database}, nil
}

func (a Adapter) Save(ctx context.Context, shipment *domain.Shipment) error {
	items := make([]ShippingItem, 0, len(shipment.Items))
	for _, item := range shipment.Items {
		items = append(items, ShippingItem{ProductCode: item.ProductCode, Quantity: item.Quantity})
	}
	model := Shipment{OrderID: shipment.OrderID, DeliveryDays: shipment.DeliveryDays, Items: items}
	res := a.db.WithContext(ctx).Create(&model)
	if res.Error == nil {
		shipment.ID = int64(model.ID)
	}
	return res.Error
}
