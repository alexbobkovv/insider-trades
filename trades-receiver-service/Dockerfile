FROM golang:1.17.6-alpine3.15 as modules

COPY go.mod go.sum /modules/
WORKDIR /modules

RUN go mod download

FROM golang:1.17.6-alpine3.15 AS builder

COPY --from=modules /go/pkg /usr/go/pkg
WORKDIR /usr/go/src/build
COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o app ./cmd/main.go

FROM scratch

COPY --from=builder /usr/go/src/build/app /app
COPY --from=builder /usr/go/src/build/.env /.env
COPY --from=builder /usr/go/src/build/config/ /config
COPY --from=builder /usr/go/src/build/migrations /migrations

CMD ["/app"]
