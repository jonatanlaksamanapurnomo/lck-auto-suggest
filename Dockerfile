FROM golang:1.21-alpine

WORKDIR /app

# Copy go mod file first
COPY go.mod ./

# Copy go sum if it exists
COPY go.sum* ./

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