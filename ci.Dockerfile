FROM alpine
ARG TARGETPLATFORM
WORKDIR /opt/alist/
COPY build/${TARGETPLATFORM//\//-}/alist ./
COPY entrypoint.sh /entrypoint.sh
RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates ffmpeg su-exec tzdata  && \
    chmod +x /opt/alist/alist /entrypoint.sh
ENV PUID=1000 PGID=1000 UMASK=022
VOLUME ["/opt/alist/data/"]
EXPOSE 5244 5245
ENTRYPOINT ["/entrypoint.sh"]
CMD ["server"]
