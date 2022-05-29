# 用于程序编译
FROM golang:1.18.2-alpine3.15 as builder
WORKDIR /build
COPY . .
RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go build -o effect .

# 用于程序运行
FROM alpine:3.15
MAINTAINER VVFock3r
WORKDIR /
COPY --from=builder /build/effect .
ENTRYPOINT ["./effect"]
CMD ["-h"]