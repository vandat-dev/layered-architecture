FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go build -mod=vendor -o app ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/app .

EXPOSE 8000

ENTRYPOINT ["./app"]
