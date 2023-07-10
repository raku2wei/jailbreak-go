FROM golang:1.20.5-alpine3.18

ARG WORKDIR=/go/app

ENV LANG=C.UTF-8 TZ=Asia/Tokyo

RUN mkdir -p $WORKDIR

WORKDIR $WORKDIR

VOLUME ["$WORKDIR"]

CMD ["/bin/ash"]
