module github.com/thiagosantosifpb/microservices/shipping

go 1.23.0

require (
	github.com/thiagosantosifpb/microservices-proto/golang/shipping v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.65.0
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.12
)

replace github.com/thiagosantosifpb/microservices-proto/golang/shipping => ../../microservices-proto/golang/shipping
