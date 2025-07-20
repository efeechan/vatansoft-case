# Dockerfile

FROM golang:1.24.5

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching deps)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the app
RUN go build -o main ./cmd/main.go

# Expose port
EXPOSE 8080

# Run it
CMD ["./main"]
