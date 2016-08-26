# Alpine linux is much smaller
FROM alpine:latest

# Install missing certificates for SSL
RUN apk add --update ca-certificates

# Push the statically-linked linux binary
ADD alpine_binary /app/alpine_binary

# Add file dependencies
ADD static /app/static
ADD templates /app/templates

# Run the web server
WORKDIR /app/
CMD ./alpine_binary
