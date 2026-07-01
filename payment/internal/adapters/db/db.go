package db

import (
	"context"
	"fmt"

	"github.com/thiagosantosifpb/microservices/payment/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	BillID     int64
	UserID     int64
	OrderID    int64
	TotalPrice float32
}

type Adapter struct{ db *gorm.DB }

func NewAdapter(dataSourceURL string) (*Adapter, error) {
	database, openErr := gorm.Open(mysql.Open(dataSourceURL), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}
	if err := database.AutoMigrate(&Payment{}); err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: database}, nil
}

func (a Adapter) Save(ctx context.Context, payment *domain.Payment) error {
	model := Payment{UserID: payment.UserID, OrderID: payment.OrderID, TotalPrice: payment.TotalPrice, BillID: payment.OrderID + 1000}
	res := a.db.WithContext(ctx).Create(&model)
	if res.Error == nil {
		payment.ID = int64(model.ID)
		payment.BillID = model.BillID
	}
	return res.Error
}
