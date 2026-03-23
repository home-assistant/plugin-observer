ARG BUILD_FROM=scratch

FROM golang:1.25-alpine3.23 AS builder

WORKDIR /workspace/observer-plugin
ARG TARGETARCH

COPY . .

# Build
RUN \
    if [ -z "${TARGETARCH}" ]; then \
        echo "TARGETARCH is not set, please use Docker BuildKit for the build." && exit 1; \
    fi \
    && case "${TARGETARCH}" in \
            amd64|arm64) ;; \
            *) echo "Unsupported TARGETARCH: ${TARGETARCH}" && exit 1 ;; \
        esac \
    && CGO_ENABLED=0 GOARCH=${TARGETARCH} go build -ldflags="-s -w" \
    && cp -f plugin-observer /workspace/observer \
    && rm -rf /workspace/observer-plugin


FROM ${BUILD_FROM}

ENV DOCKER_HOST="unix:///run/docker.sock"

WORKDIR /
COPY --from=builder /workspace/observer /usr/bin/observer
COPY rootfs /

ENTRYPOINT ["/usr/bin/observer"]

LABEL \
    io.hass.type="observer" \
    org.opencontainers.image.title="Home Assistant Observer Plugin" \
    org.opencontainers.image.description="Home Assistant Supervisor plugin monitor Supervisor" \
    org.opencontainers.image.authors="The Home Assistant Authors" \
    org.opencontainers.image.url="https://www.home-assistant.io/" \
    org.opencontainers.image.documentation="https://www.home-assistant.io/docs/" \
    org.opencontainers.image.licenses="Apache License 2.0"
