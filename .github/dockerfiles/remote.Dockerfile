FROM --platform=$BUILDPLATFORM scratch AS runner

ARG BUILDPLATFORM
ENV BUILDPLATFORM=$BUILDPLATFORM

ARG BINARY_VERSION
ENV BINARY_VERSION=$BINARY_VERSION

ARG TARGETARCH
ENV TARGETARCH=$TARGETARCH

FROM scratch AS runner

COPY --chmod=0700 --from=builder /builder/.build/api/project /project
HEALTHCHECK --interval=1m --timeout=30s --retries=3 CMD curl --fail http://localhost/status/healthcheck || exit 1
EXPOSE 80 443

LABEL org.opencontainers.image.source https://github.com/kodmain/thetiptop
ENTRYPOINT [ "/project" ]