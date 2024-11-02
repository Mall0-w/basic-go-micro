# First stage: Build the Go binary
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies (they will be cached if the files haven’t changed)
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o main .

# Second stage: Create a small final image
FROM alpine:latest

# Set up any environment variables needed
ENV PORT=8080

# Expose the port that the application will run on
EXPOSE 8080

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main /main

# Command to run the binary
CMD ["/main"]
