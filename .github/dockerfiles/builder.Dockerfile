FROM alpine:edge AS builder

RUN apk update && apk upgrade --no-cache &&\
    apk add --no-cache go go-task git alpine-sdk && \
    GOBIN=/bin go install github.com/go-task/task/v3/cmd/task@latest && \
    GOBIN=/bin go install golang.org/x/tools/cmd/goimports@latest && \
    GOBIN=/bin go install github.com/swaggo/swag/cmd/swag@latest

ARG BINARY_VERSION=build
ENV BINARY_VERSION=${BINARY_VERSION}

WORKDIR /builder
COPY . /builder

RUN task api:build

FROM scratch AS runner

COPY --chmod=0700 --from=builder /builder/.build/api/project /project
HEALTHCHECK --interval=1m --timeout=30s --retries=3 CMD curl --fail http://localhost/status/healthcheck || exit 1
EXPOSE 80 443

LABEL org.opencontainers.image.source=https://github.com/kodmain/thetiptop
ENTRYPOINT [ "/project" ]