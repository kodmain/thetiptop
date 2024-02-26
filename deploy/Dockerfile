FROM alpine:edge AS builder

ARG BINARY_VERSION
ARG BINARY_OS
ENV BINARY_VERSION=$BINARY_VERSION
ENV BINARY_OS=$BINARY_OS

LABEL org.opencontainers.image.source="https://github.com/kodmain/thetiptop"
LABEL VERSION=$BINARY_VERSION
LABEL PLATFORME=$BINARY_OS

WORKDIR /project
ADD https://github.com/kodmain/thetiptop/releases/download/$BINARY_VERSION/thetiptop-$BINARY_OS-$BINARY_ARCH /project/service
RUN chmod +x /project/service

ENTRYPOINT [ "/project/service" ]