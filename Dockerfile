# Stage 1: Build Air binary
FROM golang:1.25-alpine AS air-builder
RUN apk add --no-cache git
RUN CGO_ENABLED=0 GOOS=linux go install github.com/air-verse/air@latest

# Stage 2: Development image with Air + Go toolchain
FROM golang:1.25-alpine

# Install git (needed for go mod)
RUN apk add --no-cache git

# Copy Air binary from builder stage
COPY --from=air-builder /go/bin/air /usr/local/bin/air

WORKDIR /app

# Pre-download dependencies as a cached layer
COPY go.mod go.sum ./
RUN go mod download

# Source code is mounted as a volume at runtime for hot reload
# Air will watch for changes and rebuild automatically

EXPOSE 3000

CMD ["air", "-c", ".air.toml"]
