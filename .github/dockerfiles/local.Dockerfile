FROM scratch AS runner

ARG PROJECT_NAME
ENV PROJECT_NAME=$PROJECT_NAME

ARG REPOSITORY_NAME
ENV REPOSITORY_NAME=$REPOSITORY_NAME

WORKDIR /var/run
COPY --chmod=0777 .build/$PROJECT_NAME /var/run/$PROJECT_NAME
HEALTHCHECK --interval=1m --timeout=30s --retries=3 CMD curl --fail http://localhost/status/healthcheck || exit 1
EXPOSE 80 443

LABEL org.opencontainers.image.source https://github.com/$REPOSITORY_NAME

ENTRYPOINT [ "/var/run/$PROJECT_NAME" ]