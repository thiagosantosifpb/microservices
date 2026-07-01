package main

import (
	"github.com/thiagosantosifpb/microservices/payment/config"
	"github.com/thiagosantosifpb/microservices/payment/internal/adapters/db"
	grpcadapter "github.com/thiagosantosifpb/microservices/payment/internal/adapters/grpc"
	"github.com/thiagosantosifpb/microservices/payment/internal/application/core/api"
	"log"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("failed to connect to database. Error: %v", err)
	}
	application := api.NewApplication(dbAdapter)
	grpcAdapter := grpcadapter.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
