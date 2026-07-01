FROM golang:1.23 AS build
WORKDIR /workspace/microservices/order
COPY microservices-proto /workspace/microservices-proto
COPY microservices/order/go.mod microservices/order/go.sum* ./
RUN go mod download
COPY microservices/order/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/order ./cmd

FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/order /order
EXPOSE 3000
ENTRYPOINT ["/order"]
