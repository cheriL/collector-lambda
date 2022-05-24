FROM golang:1.15.12-alpine3.13 as builder

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go mod download && go build -ldflags="-w -s" -o /bin/collector

FROM alpine:3.13
COPY --from=builder /bin/collector .
COPY --from=builder /config ./config

ENTRYPOINT ["/collector"]