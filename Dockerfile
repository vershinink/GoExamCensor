FROM golang:1.23.1-alpine3.20 AS builder

WORKDIR /go/src/censor

COPY . .

ENV CENSOR_CONFIG_PATH=./config/config.yaml

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./censor ./cmd/main.go

FROM alpine:latest AS runner

RUN apk --no-cache add ca-certificates

WORKDIR /root

RUN mkdir -p /root/config

COPY --from=builder /go/src/censor/config ./config

COPY --from=builder /go/src/censor/censor .

EXPOSE 10503

ENTRYPOINT [ "/root/censor" ]