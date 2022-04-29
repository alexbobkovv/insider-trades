FROM golang:1.18.1-alpine3.15 as modules

COPY go.mod go.sum /modules/
WORKDIR /modules

RUN go mod download

FROM golang:1.18.1-alpine3.15 AS builder

ARG SERVICE_NAME

COPY --from=modules /go/pkg /go/pkg
WORKDIR /go/src/build/
COPY $SERVICE_NAME /go/src/build/$SERVICE_NAME
COPY go.mod go.sum ./
COPY pkg /go/src/build/pkg

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./$SERVICE_NAME/app ./$SERVICE_NAME/cmd/main.go

FROM scratch

ARG SERVICE_NAME

COPY --from=builder /go/src/build/$SERVICE_NAME/app /app
COPY --from=builder /go/src/build/$SERVICE_NAME/.env /.env
COPY --from=builder /go/src/build/$SERVICE_NAME/config/ /config
COPY --from=builder /go/src/build/$SERVICE_NAME/migrations /migrations

CMD ["/app"]
