FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY ./config/config.docker.yaml ./config/config.yaml
RUN CGO_ENABLED=0 GOOS=linux go build -o statusy ./cmd/*.go

FROM alpine
COPY --from=builder /app/statusy /statusy
COPY --from=builder /app/config /config
EXPOSE 8081
CMD ["/statusy"]
