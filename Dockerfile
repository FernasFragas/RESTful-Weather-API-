# ---------- Build Stage ----------
    FROM debian:bullseye AS builder

    # Install dependencies
    RUN apt-get update && apt-get install -y wget tar gcc libc6-dev ca-certificates
    
    # Install Go 1.23.7 manually
    ENV GOLANG_VERSION=1.23.7
    RUN wget https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
        tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz
    
    ENV PATH="/usr/local/go/bin:${PATH}"
    ENV CGO_ENABLED=1 \
        GOOS=linux \
        GOARCH=amd64
    
    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    
    RUN go build -o /app/bin/app ./cmd/web
    
    # ---------- Final Stage ----------
    FROM debian:bullseye-slim
    
    RUN apt-get update && apt-get install -y ca-certificates libsqlite3-0 && rm -rf /var/lib/apt/lists/*
    
    COPY --from=builder /app/bin/app /app/bin/app
    COPY --from=builder /app/views /app/views 
    COPY --from=builder /app/public /app/public
    # Copy bootstrap DB file (read-only copy)
    COPY --from=builder /app/bootstrap_data/weather.db /app/bootstrap_data/weather.db
 
    
    WORKDIR /app
    EXPOSE 8080
    
    CMD ["/app/bin/app"]
    
    