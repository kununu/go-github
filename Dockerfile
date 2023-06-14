FROM golang:1.20-alpine AS builder

COPY . /go/src

WORKDIR /go/src

RUN go build -o /go/bin/ghapps cmd/main.go


FROM alpine

COPY --from=builder /go/bin/ghapps /ghapps

CMD [ /ghapps ]


