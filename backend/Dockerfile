FROM golang:alpine

# Install ffmpeg and build dependencies
RUN apk add --no-cache ffmpeg build-base

WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]
