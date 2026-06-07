# Stage 1: Build Svelte UI
FROM oven/bun:1.3-alpine AS ui-builder
WORKDIR /app
COPY _ui/package.json _ui/bun.lock ./_ui/
RUN cd _ui && bun install --frozen-lockfile
COPY _ui/ ./_ui/
ENV PUBLIC_API_SERVER_ROUTE=/api
ENV PUBLIC_LOCAL_API_SERVER_ROUTE=http://localhost:8081/api
RUN cd _ui && bun run build

# Stage 2: Build Go backend
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy the built UI assets from stage 1 into the Go build context
COPY --from=ui-builder /app/_ui/build ./_ui/build
COPY ./config/config.docker.yaml ./config/config.yaml
RUN CGO_ENABLED=0 GOOS=linux go build -o statusy ./cmd/*.go

# Stage 3: Runner
FROM alpine
COPY --from=builder /app/statusy /statusy
COPY --from=builder /app/config /config
EXPOSE 8081
CMD ["/statusy"]
