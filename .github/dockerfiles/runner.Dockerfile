FROM scratch AS runner

ARG BINARY_VERSION
ENV BINARY_VERSION=$BINARY_VERSION

ARG TARGETARCH
ENV TARGETARCH=$TARGETARCH

WORKDIR /var/run
ADD --chmod=0777 https://github.com/kodmain/thetiptop/releases/download/$BINARY_VERSION/thetiptop-$TARGETARCH /var/run/project
HEALTHCHECK --interval=1m --timeout=30s --retries=3 CMD curl --fail http://localhost/v1/status/healthcheck || exit 1
EXPOSE 80 443

LABEL org.opencontainers.image.source=https://github.com/kodmain/thetiptop

ENTRYPOINT [ "/var/run/project" ]