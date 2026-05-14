ARG GO_VERSION=1.25
ARG ALPINE_VERSION=3.20

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS deps

WORKDIR /src

RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

FROM deps AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOFLAGS="-mod=readonly"

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build \
        -trimpath \
        -buildvcs=false \
        -ldflags "-s -w \
            -X main.version=${VERSION} \
            -X main.commit=${COMMIT} \
            -X main.buildDate=${BUILD_DATE}" \
        -o /out/server \
        ./cmd/server

RUN ! ldd /out/server 2>/dev/null | grep -q "=>" || (echo "binary is not static" && exit 1)

FROM gcr.io/distroless/static-debian12:nonroot AS runtime

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

# OCI labels for traceability in container registries / SBOM tooling.
LABEL org.opencontainers.image.title="Go-CollabSpace" \
      org.opencontainers.image.description="Real-time collaborative workspace backend" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.source="https://github.com/your-org/Go-CollabSpace" \
      org.opencontainers.image.licenses="MIT"

WORKDIR /app

# Copy the binary and the default config directory. In production CONFIG_PATH
# usually points at a mounted secret/configmap so this directory is optional;
# we ship it for the development MODE=development happy path.
COPY --from=builder /out/server /app/server
COPY --chown=nonroot:nonroot config /app/config

USER nonroot:nonroot

EXPOSE 8080

# Distroless has no shell, so use the exec form (no $VAR expansion possible).
ENTRYPOINT ["/app/server"]
