FROM golang:1.21.4-alpine AS builder

WORKDIR /root/gi
RUN apk add --no-cache git make gcc libc-dev

ARG GH_USER
ARG GH_TOKEN

COPY go.mod go.sum ./

ENV GH_USER=$GH_USER
ENV GH_TOKEN=$GH_TOKEN
RUN git config --global url."https://${GH_USER}:${GH_TOKEN}@github.com".insteadOf "https://github.com"

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

