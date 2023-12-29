FROM golang:1.21.4-alpine AS builder

WORKDIR /root/gi
RUN apk add --no-cache git make gcc libc-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG CGO_ENABLED 1
RUN go build cmd/main.go


FROM alpine:3.18.4 AS runner

WORKDIR /root/gi

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /root/gi/main ./gi

EXPOSE 80
ENTRYPOINT ["./gi"]

