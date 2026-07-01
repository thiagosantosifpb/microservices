FROM golang:1.23 AS build
WORKDIR /workspace/microservices/payment
COPY microservices-proto /workspace/microservices-proto
COPY microservices/payment/go.mod microservices/payment/go.sum* ./
RUN go mod download
COPY microservices/payment/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/payment ./cmd

FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/payment /payment
EXPOSE 3000
ENTRYPOINT ["/payment"]
