# ---------- Build Stage ----------
    FROM debian:bullseye AS builder

    # Install Go and required tools
    RUN apt-get update && apt-get install -y golang gcc libc6-dev ca-certificates
    
    ENV CGO_ENABLED=1 \
        GOOS=linux \
        GOARCH=amd64
    
    WORKDIR /app
    
    # Copy and download dependencies
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy the rest of the app
    COPY . .
    
    # Build the binary
    RUN go build -o /app/bin/app ./cmd/web
    
    # ---------- Final Stage ----------
    FROM debian:bullseye-slim
    
    # Install runtime dependencies only
    RUN apt-get update && apt-get install -y ca-certificates libsqlite3-0 && rm -rf /var/lib/apt/lists/*
    
    # Copy binary and static files
    COPY --from=builder /app/bin/app /app/bin/app
    COPY --from=builder /app/views /app/views 
    COPY --from=builder /app/public /app/public 
    
    # Set working directory and expose port
    WORKDIR /app
    EXPOSE 8080
    
    CMD ["/app/bin/app"]
    