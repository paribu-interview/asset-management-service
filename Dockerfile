# Step 1: Build the Go application
FROM golang:1.23 AS builder

WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ams ./cmd/main.go

# Verify the binary
RUN ls -l /app

# Step 2: Create a lightweight runtime image
FROM alpine:3.18

RUN apk --no-cache add ca-certificates bash busybox-extras

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/ams /app/ams
COPY --from=builder /app/envs /app/envs
COPY --from=builder /app/db/migrations /app/db/migrations
COPY wait-for-it.sh /app/wait-for-it.sh

# Verify the binary in the runtime image
RUN ls -l /app

EXPOSE 8081

CMD ["/app/wait-for-it.sh", "postgres", "5432", "/app/ams"]
