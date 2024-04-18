FROM alpine:edge AS builder

RUN echo "hello" > /var/run/$(echo | awk '{srand(); print int(rand()*100) + 1}')

FROM scratch AS runner

WORKDIR /var/run
COPY --from=builder /var/run /var/run
HEALTHCHECK --interval=1m --timeout=30s --retries=3 CMD curl --fail http://localhost/status/healthcheck || exit 1
EXPOSE 80 443

ENTRYPOINT [ "/var/run/project" ]