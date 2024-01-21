FROM node:18-alpine as frontend

WORKDIR /frontend_dir
COPY ./_ui/package.json .
RUN npm install
COPY ./_ui .
RUN npm run build


FROM golang:1.21 as backend

WORKDIR /backend_dir
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN rm -rf ./_ui
RUN mkdir -p ./_ui
COPY --from=frontend /frontend_dir/build ./_ui/build
RUN CGO_ENABLED=0 GOOS=linux go build -o /backend_dir/statusy cmd/main/*



FROM alpine:3.18

RUN apk add --no-cache ca-certificates tzdata
COPY --from=backend /backend_dir/statusy /app/statusy
EXPOSE 8080

CMD ["/app/statusy"]