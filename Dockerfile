FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY ./config ./
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o statusy ./cmd/*.go

FROM scratch
COPY --from=builder /app/statusy /statusy
CMD ["/statusy"]
