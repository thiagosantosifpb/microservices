# Instruções passo a passo para executar

## 1. Pré-requisitos

- Docker e Docker Compose instalados;
- Go 1.23 ou superior, caso deseje executar sem Docker;
- `grpcurl`, caso deseje testar chamadas gRPC pelo terminal.

## 2. Executar com Docker Compose

A partir da pasta `microservices`:

```bash
docker compose up --build
```

Aguarde a criação do banco MySQL e a inicialização dos três microsserviços.

## 3. Portas

- Order: `localhost:3000`;
- Payment: `localhost:3001`;
- Shipping: `localhost:3002`;
- MySQL: `localhost:3306`.

## 4. Teste de sucesso

```bash
grpcurl -plaintext \
  -d '{"costumer_id":123,"order_items":[{"product_code":"NOTEBOOK","quantity":2,"unit_price":100},{"product_code":"MOUSE","quantity":3,"unit_price":50}],"total_price":350}' \
  localhost:3000 Order/Create
```

## 5. Teste de erro por produto inexistente

```bash
grpcurl -plaintext \
  -d '{"costumer_id":123,"order_items":[{"product_code":"PRODUTO_INVALIDO","quantity":1,"unit_price":10}],"total_price":10}' \
  localhost:3000 Order/Create
```

## 6. Teste de erro por mais de 50 itens

```bash
grpcurl -plaintext \
  -d '{"costumer_id":123,"order_items":[{"product_code":"MOUSE","quantity":51,"unit_price":10}],"total_price":510}' \
  localhost:3000 Order/Create
```

## 7. Teste de erro por pagamento acima de 1000

```bash
grpcurl -plaintext \
  -d '{"costumer_id":123,"order_items":[{"product_code":"NOTEBOOK","quantity":5,"unit_price":250}],"total_price":1250}' \
  localhost:3000 Order/Create
```

## 8. Executar sem Docker

Suba um MySQL com as bases `order`, `payment` e `shipping`, depois execute cada serviço em terminais separados:

```bash
cd microservices/payment
DATA_SOURCE_URL='root:minhasenha@tcp(127.0.0.1:3306)/payment?charset=utf8mb4&parseTime=True&loc=Local' APPLICATION_PORT=3001 ENV=development go run ./cmd
```

```bash
cd microservices/shipping
DATA_SOURCE_URL='root:minhasenha@tcp(127.0.0.1:3306)/shipping?charset=utf8mb4&parseTime=True&loc=Local' APPLICATION_PORT=3002 ENV=development go run ./cmd
```

```bash
cd microservices/order
DATA_SOURCE_URL='root:minhasenha@tcp(127.0.0.1:3306)/order?charset=utf8mb4&parseTime=True&loc=Local' APPLICATION_PORT=3000 ENV=development PAYMENT_SERVICE_URL=localhost:3001 SHIPPING_SERVICE_URL=localhost:3002 go run ./cmd
```
