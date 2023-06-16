# Builder image
FROM golang:1.20-alpine AS builder

# Copy source code to container
COPY . /go/src

WORKDIR /go/src

# Build source code
RUN go build -o /go/bin/ghapps cmd/main.go

# Final image
FROM scratch

# Copy SSL certs from buider. This is needed for the app to run
COPY --from=builder /etc/ssl /etc/ssl

# Copy the binary from builder
COPY --from=builder /go/bin/ghapps /ghapps

ENTRYPOINT [ "/ghapps" ]


