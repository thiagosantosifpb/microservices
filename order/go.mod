module github.com/thiagosantosifpb/microservices/order

go 1.23.0

require (
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/thiagosantosifpb/microservices-proto/golang/order v0.0.0-00010101000000-000000000000
	github.com/thiagosantosifpb/microservices-proto/golang/payment v0.0.0-00010101000000-000000000000
	github.com/thiagosantosifpb/microservices-proto/golang/shipping v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.65.0
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.12
)

replace github.com/thiagosantosifpb/microservices-proto/golang/order => ../../microservices-proto/golang/order
replace github.com/thiagosantosifpb/microservices-proto/golang/payment => ../../microservices-proto/golang/payment
replace github.com/thiagosantosifpb/microservices-proto/golang/shipping => ../../microservices-proto/golang/shipping
