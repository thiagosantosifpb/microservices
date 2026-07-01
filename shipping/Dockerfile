FROM golang:1.23 AS build
WORKDIR /workspace/microservices/shipping
COPY microservices-proto /workspace/microservices-proto
COPY microservices/shipping/go.mod microservices/shipping/go.sum* ./
RUN go mod download
COPY microservices/shipping/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/shipping ./cmd

FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/shipping /shipping
EXPOSE 3000
ENTRYPOINT ["/shipping"]
