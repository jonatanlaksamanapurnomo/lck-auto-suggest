FROM golang:1.21-alpine

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Create data directory and copy TSV
RUN mkdir -p /app/data
COPY data/cities_canada-usa.tsv /app/data/

# Build the application
RUN go build -o main ./cmd/api/main.go

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./main"]