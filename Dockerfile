FROM golang:1.22.1-alpine3.19

ARG WORKDIR=/go/app

ENV LANG=C.UTF-8 TZ=Asia/Tokyo

RUN mkdir -p $WORKDIR

WORKDIR $WORKDIR

VOLUME ["$WORKDIR"]

CMD ["/bin/ash"]
