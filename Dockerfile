# go.mod が go 1.22 を要求するため(Ebitengine v2.8 対応)
FROM golang:1.22-alpine3.20

ARG WORKDIR=/go/app

ENV LANG=C.UTF-8 TZ=Asia/Tokyo

RUN mkdir -p $WORKDIR

WORKDIR $WORKDIR

VOLUME ["$WORKDIR"]

CMD ["/bin/ash"]
