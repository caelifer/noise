# Milti-stage build

## Builder stage
FROM golang:latest as builder

# Maintener info
LABEL maintainer="Timour Ezeev <timour.ezeev@me.com>"

# Switch to Workspace directory in the builder
WORKDIR /build

# Copy content of current directory to workspace in the container
COPY . .

# Download all Go dependencies

# If proxy required for connecting to the Internet, set these:
# ARG https_proxy=http://http-proxy:port
# ARG http_proxy=http://http-proxy:port
RUN go mod tidy

# Build our executable
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o noise .

## Rutime image
FROM alpine:latest

# If proxy required for connecting to the Internet, set these:
# ARG https_proxy=http://http-proxy:port
# ARG http_proxy=http://http-proxy:port
RUN apk --no-cache add ca-certificates

# Copy prebuild binary to the application folder
WORKDIR /app
COPY --from=builder /build/noise .

# Add non-priveleged user to run our app
RUN adduser appuser -D -H
RUN chown -R appuser /app
USER appuser

# Expose ports
EXPOSE 8080

# Set required environment variable
ARG APP_HTTP_PORT=8080

# Execute our binary
CMD ["/app/noise"]
