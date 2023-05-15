FROM golang:1.20.3 AS build
COPY . /opt
WORKDIR /opt
ENV GOPROXY https://goproxy.cn
ENV GO111MODULE on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o app ./src

FROM alpine
COPY --from=build /opt/app /opt/app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && apk add tzdata
ENTRYPOINT ["/opt/app"]