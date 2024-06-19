# Step 1: Build the Go binary
FROM golang:1.21.5 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd

# Step 2: Build a small image for the Go binary
FROM scratch

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main /app/main
COPY --from=builder /app/ui/index.html /ui/index.html

# Expose port 5001 to the outside world
EXPOSE 5001

# Command to run the executable
ENTRYPOINT ["/app/main"]
