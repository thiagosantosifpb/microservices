package domain

import "time"

type ShippingItem struct {
	ProductCode string `json:"product_code"`
	Quantity    int32  `json:"quantity"`
}

type Shipment struct {
	ID           int64          `json:"id"`
	OrderID      int64          `json:"order_id"`
	Items        []ShippingItem `json:"items"`
	DeliveryDays int32          `json:"delivery_days"`
	CreatedAt    int64          `json:"created_at"`
}

func NewShipment(orderID int64, items []ShippingItem) Shipment {
	return Shipment{OrderID: orderID, Items: items, CreatedAt: time.Now().Unix()}
}

func (s *Shipment) TotalQuantity() int32 {
	var total int32
	for _, item := range s.Items {
		total += item.Quantity
	}
	return total
}

func (s *Shipment) CalculateDeliveryDays() int32 {
	return 1 + (s.TotalQuantity() / 5)
}
