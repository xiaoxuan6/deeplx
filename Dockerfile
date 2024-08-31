FROM golang:1.22.5-alpine3.20 AS build-dev

WORKDIR /go/app/src

COPY . .

RUN apk add --no-cache upx tzdata || \
    go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod tidy && \
    go build -ldflags="-s -w" -o deeplx ./api/main.go && \
    [ -e /usr/bin/upx ] && upx deeplx || echo

FROM alpine

COPY --from=build-dev /go/app/src/deeplx .
COPY --from=build-dev /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-dev /usr/share/zoneinfo /usr/share/zoneinfo
COPY --link blacklist.txt .

ENV ROUTER_PATH=""
ENV TZ=Asia/Shanghai

EXPOSE 8311

ENTRYPOINT ["./deeplx"]