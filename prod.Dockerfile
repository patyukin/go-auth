FROM golang:1.22.3-alpine3.19 AS builder

RUN apk add --no-cache gcc musl-dev sqlite

ENV CGO_ENABLED=1

COPY . /app
WORKDIR /app

RUN go mod download
RUN go mod tidy
RUN go build -o ./bin/auth cmd/main.go

FROM alpine:3.19

RUN apk add --no-cache sqlite

WORKDIR /app
COPY --from=builder /app/bin/auth .

CMD ["./auth", "--config", "config.yaml"]
