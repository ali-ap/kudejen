# Stage 1: Build the executable
FROM golang:1.22.0 AS builder
LABEL authors="Ali"
# Set necessary environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOLANG_VERSION=1.22.0

# Create a directory for the app
WORKDIR /app

# Copy and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the entire source code
COPY . .

# Set a default version number
ARG DEFAULT_VERSION="1.0.0"

# Retrieve version number from environment variable or use default
ARG APP_VERSION
ENV VERSION=${APP_VERSION:-$DEFAULT_VERSION}
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X 'main.version=${VERSION}'" -o kudejen src/cmd/main.go

# Stage 2: Generate a minimal image with only the binary
FROM alpine:latest
RUN apk --no-cache add ca-certificates


# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the pre-built binary from the previous stage
WORKDIR /home/appuser/
COPY --from=builder /app/kudejen /home/appuser/
COPY --from=builder app/src/internal/config/config.yml /home/appuser//src/internal/config/config.yml
COPY --from=builder app/src/internal/config/kube.config /home/appuser//src/internal/config/kube.config

# Change ownership of the copied files to the non-root user
RUN chown -R appuser:appgroup /home/appuser

# Switch to the non-root user
USER appuser

# Expose port for gRPC
EXPOSE 8080
EXPOSE 8081

# Command to run the executable
CMD ["./kudejen"]
