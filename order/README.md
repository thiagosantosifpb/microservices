# Order Microservice

Microsserviço `Order` desenvolvido em Go com arquitetura hexagonal, gRPC e persistência em MySQL usando GORM.

## Estrutura

```text
order/
├── cmd/
│   └── main.go
├── config/
│   └── config.go
├── internal/
│   ├── adapters/
│   │   ├── db/
│   │   │   └── db.go
│   │   └── grpc/
│   │       └── server.go
│   ├── application/
│   │   └── core/
│   │       ├── api/
│   │       │   └── api.go
│   │       └── domain/
│   │           └── order.go
│   └── ports/
│       ├── api.go
│       └── db.go
├── docker-compose.yml
├── go.mod
└── .gitignore
```

## Pré-requisitos

- Go 1.23 ou superior
- Docker
- grpcurl, para testar chamadas gRPC

## Importante sobre a estrutura de pastas

Para o `replace` do `go.mod` funcionar sem alterações, mantenha os dois repositórios lado a lado:

```text
pasta-de-trabalho/
├── microservices-proto/
└── microservices/
    └── order/
```

## Instalar dependências

```bash
cd microservices/order
go mod tidy
```

## Subir o MySQL

Opção 1, usando Docker Compose:

```bash
cd microservices/order
docker compose up -d
```

Opção 2, usando o comando do enunciado:

```bash
docker run -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=minhasenha \
  -e MYSQL_DATABASE=order \
  mysql:8.4
```

## Executar o microsserviço

```bash
cd microservices/order
DATA_SOURCE_URL='root:minhasenha@tcp(127.0.0.1:3306)/order?charset=utf8mb4&parseTime=True&loc=Local' \
APPLICATION_PORT=3000 \
ENV=development \
go run ./cmd
```

## Testar com grpcurl

O campo foi mantido como `costumer_id`, exatamente como aparece no `order.proto` da prática.

```bash
grpcurl -plaintext \
  -d '{"costumer_id":123,"order_items":[{"product_code":"prod","quantity":4,"unit_price":12}],"total_price":48}' \
  localhost:3000 Order/Create
```

Resposta esperada:

```json
{
  "orderId": 1
}
```

por:

```text
github.com/SEU_USUARIO
```
