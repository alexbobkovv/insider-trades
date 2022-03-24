FROM golang:1.17.6-alpine3.15 AS builder

WORKDIR /usr/go/src/build
COPY . .
RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

FROM scratch

COPY --from=builder /usr/go/src/build/app /
COPY --from=builder /usr/go/src/build/config/config.yml /config
COPY --from=builder /usr/go/src/build/migrations /migrations

CMD ["/app"]
