package main

import (
	"github.com/thiagosantosifpb/microservices/shipping/config"
	"github.com/thiagosantosifpb/microservices/shipping/internal/adapters/db"
	grpcadapter "github.com/thiagosantosifpb/microservices/shipping/internal/adapters/grpc"
	"github.com/thiagosantosifpb/microservices/shipping/internal/application/core/api"
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
