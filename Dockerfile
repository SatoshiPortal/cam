FROM golang:latest as builder

ADD . /src

WORKDIR /src
RUN ./build.sh

FROM scratch

COPY --from=builder /src/cam /cam

WORKDIR /data

ENTRYPOINT ["/cam"]