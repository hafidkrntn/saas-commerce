# ---------- BUILDER ----------
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./cmd/main.go

# ---------- FINAL IMAGE ----------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
ENV TZ=Asia/Jakarta

COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations

# Security: don't run as root
RUN adduser -D runner
USER runner

EXPOSE 8000
CMD ["./app"]
