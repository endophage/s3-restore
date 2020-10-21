FROM golang:1.15-alpine as builder

RUN apk add --no-cache git make curl openssl

ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOOS=linux \
    COARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build . \
    && cp s3-restore /usr/local/bin/

FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/* /usr/local/bin

RUN adduser -D s3-restore
USER s3-restore

ENTRYPOINT [ "/usr/local/bin/s3-restore" ]
