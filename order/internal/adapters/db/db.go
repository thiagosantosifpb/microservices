package db

import (
	"fmt"

	"github.com/thiagosantosifpb/microservices/order/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code        string `gorm:"uniqueIndex;size:80;not null"`
	Description string
}

type Order struct {
	gorm.Model
	CustomerID   int64
	Status       string
	DeliveryDays int32
	OrderItems   []OrderItem
}

type OrderItem struct {
	gorm.Model
	ProductCode string
	UnitPrice   float32
	Quantity    int32
	OrderID     uint
}

type Adapter struct{ db *gorm.DB }

func NewAdapter(dataSourceURL string) (*Adapter, error) {
	database, openErr := gorm.Open(mysql.Open(dataSourceURL), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}
	if err := database.AutoMigrate(&Product{}, &Order{}, &OrderItem{}); err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	adapter := &Adapter{db: database}
	return adapter, adapter.seedProducts()
}

func (a Adapter) seedProducts() error {
	products := []Product{{Code: "NOTEBOOK", Description: "Notebook"}, {Code: "MOUSE", Description: "Mouse"}, {Code: "KEYBOARD", Description: "Teclado"}, {Code: "MONITOR", Description: "Monitor"}, {Code: "HEADSET", Description: "Headset"}}
	for _, p := range products {
		if err := a.db.FirstOrCreate(&p, Product{Code: p.Code}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (a Adapter) ProductExists(productCode string) (bool, error) {
	var count int64
	if err := a.db.Model(&Product{}).Where("code = ?", productCode).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (a Adapter) Get(id string) (domain.Order, error) {
	var orderEntity Order
	res := a.db.Preload("OrderItems").First(&orderEntity, id)
	var items []domain.OrderItem
	for _, it := range orderEntity.OrderItems {
		items = append(items, domain.OrderItem{ProductCode: it.ProductCode, UnitPrice: it.UnitPrice, Quantity: it.Quantity})
	}
	return domain.Order{ID: int64(orderEntity.ID), CustomerID: orderEntity.CustomerID, Status: orderEntity.Status, DeliveryDays: orderEntity.DeliveryDays, OrderItems: items, CreatedAt: orderEntity.CreatedAt.Unix()}, res.Error
}

func (a Adapter) Save(order *domain.Order) error {
	items := make([]OrderItem, 0, len(order.OrderItems))
	for _, it := range order.OrderItems {
		items = append(items, OrderItem{ProductCode: it.ProductCode, UnitPrice: it.UnitPrice, Quantity: it.Quantity})
	}
	model := Order{CustomerID: order.CustomerID, Status: order.Status, DeliveryDays: order.DeliveryDays, OrderItems: items}
	res := a.db.Create(&model)
	if res.Error == nil {
		order.ID = int64(model.ID)
	}
	return res.Error
}

func (a Adapter) UpdateStatus(orderID int64, status string) error {
	return a.db.Model(&Order{}).Where("id = ?", orderID).Update("status", status).Error
}
func (a Adapter) UpdateDeliveryDays(orderID int64, deliveryDays int32) error {
	return a.db.Model(&Order{}).Where("id = ?", orderID).Update("delivery_days", deliveryDays).Error
}
