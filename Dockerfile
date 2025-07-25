# -------- Build Stage --------
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

# -------- Final Stage --------
FROM alpine:latest

# Install Docker CLI
RUN apk add --no-cache docker-cli

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8001

CMD ["./app"]
