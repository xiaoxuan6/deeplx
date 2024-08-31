FROM golang:1.22.5-alpine3.20 AS build-dev

WORKDIR /go/app/src

COPY . .

RUN apk add --no-cache upx || \
    go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod tidy && \
    go build -ldflags="-s -w" -o deeplx ./api/main.go && \
    [ -e /usr/bin/upx ] && upx deeplx || echo

FROM scratch

COPY --from=build-dev /go/app/src/deeplx .
ENV ROUTER_PATH=""

ENTRYPOINT ["/deeplx"]