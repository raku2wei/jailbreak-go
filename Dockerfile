FROM golang:1.18

ARG WORKDIR=/go/app

ENV LANG=C.UTF-8 TZ=Asia/Tokyo

RUN mkdir -p $WORKDIR

WORKDIR $WORKDIR

VOLUME ["$WORKDIR"]

CMD ["/bin/bash"]
