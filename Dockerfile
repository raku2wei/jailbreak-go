# go.mod の go / toolchain 行(Go 1.27)に合わせる。
# 公式 golang イメージは GOTOOLCHAIN=local のため、go.mod より古い Go のイメージではビルドできない。
FROM golang:1.27-alpine3.24

ARG WORKDIR=/go/app

ENV LANG=C.UTF-8 TZ=Asia/Tokyo

RUN mkdir -p $WORKDIR

WORKDIR $WORKDIR

VOLUME ["$WORKDIR"]

CMD ["/bin/ash"]
