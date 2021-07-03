FROM alpine:latest
RUN apk --update add ca-certificates \
                     mailcap \
                     curl

HEALTHCHECK --start-period=2s --interval=5s --timeout=3s \
  CMD curl -f http://localhost/health || exit 1

VOLUME /srv
EXPOSE 80

RUN mkdir /var/lib/filebrowser && chmod a+w /var/lib/filebrowser
COPY .docker.json /.filebrowser.json
COPY filebrowser /filebrowser

ENTRYPOINT [ "/filebrowser" ]