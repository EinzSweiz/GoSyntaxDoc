# ✅ Stage 1: Build the Go application
FROM golang:1.24 AS builder

WORKDIR /app

# ✅ Copy go.mod and go.sum first (cache dependencies)
COPY go.mod go.sum ./
RUN go mod download

# ✅ Copy the application source code
COPY . .

# ✅ Build the Go binary (Force static binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o fiber-app cmd/main.go

# ✅ Stage 2: Create a minimal image with the built binary
FROM alpine:latest

WORKDIR /root/

# ✅ Install CA certificates (needed for HTTPS requests)
RUN apk --no-cache add ca-certificates

# ✅ Copy the compiled binary from the builder stage
COPY --from=builder /app/fiber-app /root/fiber-app

# ✅ Ensure the binary is executable
RUN chmod +x /root/fiber-app

# ✅ Expose the Fiber app port
EXPOSE 3002

# ✅ Run the application
CMD ["/root/fiber-app"]
