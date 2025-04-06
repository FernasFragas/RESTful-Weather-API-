# Use the official Golang image as a build stage
FROM golang:1.23.7 as builder

# Disable CGO for a static binary (no glibc required)
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary
RUN go build -o /app/bin/app ./cmd/web

# ---------- Final Stage ----------
FROM gcr.io/distroless/static:nonroot


# Copy only the binary
COPY --from=builder /app/bin/app /app/bin/app
COPY --from=builder /app/views /app/views 
COPY --from=builder /app/public /app/public 

# Set working directory and expose port
WORKDIR /app
EXPOSE 8080

CMD ["/app/bin/app"]