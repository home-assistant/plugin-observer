ARG BUILD_FROM

FROM golang:1.18-alpine3.16 AS builder

WORKDIR /workspace/observer-plugin
ARG BUILD_ARCH

COPY . .

# Build
RUN \
        if [ "${BUILD_ARCH}" = "armhf" ]; then \
            CGO_ENABLED=0 GOARM=6 GOARCH=arm go build -ldflags="-s -w"; \
        elif [ "${BUILD_ARCH}" = "armv7" ]; then \
            CGO_ENABLED=0 GOARM=7 GOARCH=arm go build -ldflags="-s -w"; \
        elif [ "${BUILD_ARCH}" = "aarch64" ]; then \
            CGO_ENABLED=0 GOARCH=arm64 go build -ldflags="-s -w"; \
        elif [ "${BUILD_ARCH}" = "i386" ]; then \
            CGO_ENABLED=0 GOARCH=386 go build -ldflags="-s -w"; \
        elif [ "${BUILD_ARCH}" = "amd64" ]; then \
            CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w"; \
        else \
            exit 1; \
        fi \
    && cp -f plugin-observer /workspace/observer \
    && rm -rf /workspace/observer-plugin


FROM ${BUILD_FROM}

WORKDIR /
COPY --from=builder /workspace/observer /usr/bin/observer
COPY rootfs /

ENTRYPOINT ["/usr/bin/observer"]