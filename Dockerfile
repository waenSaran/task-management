FROM golang:1.23.4

WORKDIR /app

COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code from current directory to working directory in the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# run command through the binary file (built file)
CMD ["./main"]