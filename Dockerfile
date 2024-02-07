FROM rss3/go-builder AS base

WORKDIR /root/gi
RUN apk add --no-cache git make gcc libc-dev

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .

FROM base AS builder

ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go build cmd/main.go

FROM rss3/go-runtime AS runner

WORKDIR /root/gi

COPY --from=builder /root/gi/main ./gi

EXPOSE 80
ENTRYPOINT ["./gi"]

