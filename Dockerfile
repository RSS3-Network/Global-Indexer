FROM ghcr.io/rss3-network/go-image/go-builder AS base

WORKDIR /root/gi

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .

# Download GeoLite2-City.mmdb from Henry's Google Drive
RUN mkdir -p common/geolite2/mmdb && touch common/geolite2/mmdb/.geoipupdate.lock
RUN curl -L "https://drive.google.com/uc?export=download&id=1xyJc_0rupY5MZzdCTOr7sDL--j7r5i0w" -o common/geolite2/mmdb/GeoLite2-City.mmdb

FROM base AS builder

ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go build cmd/main.go

FROM ghcr.io/rss3-network/go-image/go-runtime AS runner

WORKDIR /root/gi

COPY --from=builder /root/gi/main ./gi

EXPOSE 80
ENTRYPOINT ["./gi"]
