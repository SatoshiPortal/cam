FROM golang:latest as builder

ADD . /src

WORKDIR /src
RUN ./build.sh

FROM scratch

COPY --from=builder /src/cam /cam
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

WORKDIR /data

ENTRYPOINT ["/cam"]