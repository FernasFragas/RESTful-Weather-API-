   # Use the official Golang image as a build stage
   FROM golang:1.23.7 as builder

   # Set the working directory inside the container
   WORKDIR /app

   # Copy go.mod and go.sum files
   COPY go.mod go.sum ./

   # Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
   RUN go mod download

   # Copy the source code into the container
   COPY . .

   # Build the Go app
   RUN go build -o /app/bin/app ./cmd/web

   # Use a more recent Debian image for the final stage
   FROM debian:bullseye-slim

   # Install necessary packages
   RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

   # Copy the binary from the builder stage
   COPY --from=builder /app/bin/app /app/bin/app

   # Expose the port the app runs on
   EXPOSE 8080

   # Command to run the executable
   CMD ["/app/bin/app"]