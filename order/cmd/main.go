package main

import (
	"log"

	"github.com/thiagosantosifpb/microservices/order/config"
	"github.com/thiagosantosifpb/microservices/order/internal/adapters/db"
	grpcadapter "github.com/thiagosantosifpb/microservices/order/internal/adapters/grpc"
	paymentadapter "github.com/thiagosantosifpb/microservices/order/internal/adapters/payment"
	shippingadapter "github.com/thiagosantosifpb/microservices/order/internal/adapters/shipping"
	"github.com/thiagosantosifpb/microservices/order/internal/application/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("failed to connect to database. Error: %v", err)
	}
	paymentAdapter, err := paymentadapter.NewAdapter(config.GetPaymentServiceURL())
	if err != nil {
		log.Fatalf("failed to initialize payment stub. Error: %v", err)
	}
	shippingAdapter, err := shippingadapter.NewAdapter(config.GetShippingServiceURL())
	if err != nil {
		log.Fatalf("failed to initialize shipping stub. Error: %v", err)
	}
	application := api.NewApplication(dbAdapter, paymentAdapter, shippingAdapter)
	grpcAdapter := grpcadapter.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
