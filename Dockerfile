# Builder image
FROM golang:1.24-alpine AS builder

# Copy source code to container
COPY . /go/src

WORKDIR /go/src

# Build source code
RUN go build -o /go/bin/ghapps cmd/main.go

# Final image
FROM alpine

# Copy the binary from builder
COPY --from=builder /go/bin/ghapps /ghapps

ENTRYPOINT [ "/ghapps" ]


