FROM golang:1.18.1 as builder
# Define build env
ENV GOOS linux
# Add a work directory
WORKDIR /app
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy app files
COPY . .
# Build app
RUN go build -o app

FROM debian:buster
# # Copy built binary from builder
COPY --from=builder app .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Expose port
EXPOSE 8080
# Exec built binary
CMD ./app