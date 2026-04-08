FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# Deps first — cached unless go.mod/go.sum change
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o mcp_server \
    ./cmd/mcp_server

# ---- runtime ----
FROM alpine:3.21

# ca-certificates needed for outbound HTTPS (dashboard API)
RUN apk add --no-cache ca-certificates && \
    update-ca-certificates

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/mcp_server .

ENV GIN_MODE=release

USER app

ENTRYPOINT ["/app/mcp_server"]
