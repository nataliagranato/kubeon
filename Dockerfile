FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o kubeon .

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /app/kubeon .

USER 1000

ENTRYPOINT ["./kubeon"]
