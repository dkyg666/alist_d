FROM golang:alpine AS builder
ARG VERSION=dev
WORKDIR /build
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash gcc git musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN bash build.sh $VERSION

FROM alpine
WORKDIR /opt/alist/
COPY --from=builder /build/alist ./
COPY entrypoint.sh /entrypoint.sh
RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates ffmpeg su-exec tzdata && \
    chmod +x /opt/alist/alist /entrypoint.sh
ENV PUID=0 PGID=0 UMASK=022
VOLUME /opt/alist/data/
EXPOSE 5244 5245
ENTRYPOINT ["/entrypoint.sh"]
CMD ["server"]
