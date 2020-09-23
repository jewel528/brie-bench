FROM golang:1.14-alpine
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct
ENV ARTIFICATS=/artifacts
WORKDIR /brie
COPY . .

RUN apk add --no-cache \
    make \
    git \
    bash \
    curl \
    gcc \
    g++

RUN mkdir $ARTIFICATS
RUN ls /
CMD ["./run.sh"]