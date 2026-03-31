FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

FROM alpine:3.22

WORKDIR /app

RUN adduser -D appuser
COPY --from=builder /app/bin/api /app/api

USER appuser

EXPOSE 8080

CMD ["/app/api"]