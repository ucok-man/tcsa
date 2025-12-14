# Build stage
FROM golang:1.25.0-alpine AS builder

RUN apk add --no-cache make
WORKDIR /app

COPY go.mod go.sum ./
COPY vendor/ ./vendor/
COPY . .
RUN make build

# Runtime stage
FROM alpine:latest

# Create non-root user
RUN addgroup -g 1000 tsca.api && \
    adduser -D -u 1000 -G tsca.api tsca.api

WORKDIR /app

COPY --from=builder /app/bin/api .
COPY --from=builder /app/cmd/api/docs ./cmd/api/docs

RUN chown -R tsca.api:tsca.api /app
USER tsca.api
EXPOSE 4000

CMD ["./api"]