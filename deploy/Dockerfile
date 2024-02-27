FROM --platform=$BUILDPLATFORM alpine:edge AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM

ARG BINARY_VERSION
ARG BINARY_OS
ENV BINARY_VERSION=$BINARY_VERSION
ENV BINARY_OS=$BINARY_OS

RUN if [ -n "$TARGETPLATFORM" ]; then \
        PLATFORM=$(echo "$TARGETPLATFORM" | cut -f2 -d '/'); \
    else \
        PLATFORM=$(echo "$BUILDPLATFORM" | cut -f2 -d '/'); \
    fi && export BINARY_ARCH="$PLATFORM";

LABEL org.opencontainers.image.source="https://github.com/kodmain/thetiptop"
LABEL VERSION=$BINARY_VERSION

WORKDIR /project
ADD https://github.com/kodmain/thetiptop/releases/download/$BINARY_VERSION/thetiptop-linux-$BINARY_ARCH /project/service
RUN chmod +x /project/service

ENTRYPOINT [ "/project/service" ]