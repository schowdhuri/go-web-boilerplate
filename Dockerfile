# Build stage
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache nodejs npm
RUN npm install -g pnpm
RUN pnpm i --frozen-lockfile
RUN make build

# Production stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin ./bin
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/internal/templates ./internal/templates

EXPOSE 8080
CMD ["./bin/main"]
