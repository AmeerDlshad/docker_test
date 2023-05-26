FROM alpine:edge as builder
# Define build env
RUN apk add --no-cache --update go gcc g++ 
# Add a work directory
WORKDIR /app
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy app files
COPY . .
# Build app

RUN CGO_ENABLED=1 GOOS=linux go build -o app

FROM alpine:3.14
# # Copy built binary from builder
COPY --from=builder app .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN apk add --no-cache --update curl

# Expose port
EXPOSE 8080
# Exec built binary
CMD ./app