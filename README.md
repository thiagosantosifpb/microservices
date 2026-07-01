# Microservices gRPC

Projeto completo da prática de microsserviços com gRPC em Go, organizado em arquitetura hexagonal.

## Microsserviços

- `order`: recebe pedidos, valida estoque, limita a quantidade total a 50 itens, registra o pedido, chama `payment` e, em caso de pagamento aprovado, chama `shipping`.
- `payment`: registra cobranças e recusa pagamentos acima de R$ 1.000,00 com erro gRPC `InvalidArgument`.
- `shipping`: calcula o prazo de entrega a partir da quantidade total de unidades. O prazo mínimo é 1 dia e, a cada 5 unidades, soma-se 1 dia.

## Executando com Docker Compose

A partir da pasta `microservices`:

```bash
docker compose up --build
```

O serviço `order` ficará em `localhost:3000`, `payment` em `localhost:3001` e `shipping` em `localhost:3002`.

## Teste com grpcurl

```bash
grpcurl -plaintext \
  -d '{"costumer_id":123,"order_items":[{"product_code":"NOTEBOOK","quantity":2,"unit_price":100},{"product_code":"MOUSE","quantity":3,"unit_price":50}],"total_price":350}' \
  localhost:3000 Order/Create
```

Produtos cadastrados automaticamente no banco do `order`: `NOTEBOOK`, `MOUSE`, `KEYBOARD`, `MONITOR`, `HEADSET`.
