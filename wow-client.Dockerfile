FROM golang:1.22 AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o client ./cmd/client

FROM scratch

COPY --from=builder /build/client /

ENTRYPOINT ["/client"]
