FROM alpine:edge AS builder

ARG BINARY_VERSION
ARG BINARY_OS
ARG BINARY_ARCH
ENV BINARY_VERSION=$BINARY_VERSION
ENV BINARY_OS=$BINARY_OS
ENV BINARY_ARCH=$BINARY_ARCH

WORKDIR /project
ADD https://github.com/kodmain/thetiptop/releases/download/$BINARY_VERSION/thetiptop-$BINARY_OS-$BINARY_ARCH /project/service
RUN chmod +x /project/service

ENTRYPOINT [ "/project/service" ]