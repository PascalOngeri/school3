# Stage 1: Build stage
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the entire application source code and additional folders into the container
COPY . .

# Build the Go app as a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o feego .

# Stage 2: Production stage
FROM debian:bullseye-slim

# Install necessary runtime libraries
RUN apt-get update && \
    apt-get install -y --no-install-recommends libc6 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/feego .

# Copy templates, static, assets, includes, and uploads folders
COPY --from=builder /app/templates /root/templates
COPY --from=builder /app/assets /root/assets
COPY --from=builder /app/includes /root/includes
COPY --from=builder /app/uploads /root/uploads

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./feego"]
