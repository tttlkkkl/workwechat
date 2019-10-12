FROM golang:1.13.1-alpine3.10 as build_deps
RUN apk add --no-cache git
WORKDIR /workspace
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io
COPY go.mod .
COPY go.sum .
RUN go mod download
FROM build_deps AS build

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o app -work /workspace/dev
FROM alpine:3.10
LABEL MAINTAINER="tttlkkkl <tttlkkkl@aliyun.com>"
ENV TZ "Asia/Shanghai"
ENV TERM xterm

RUN echo 'https://mirrors.aliyun.com/alpine/v3.10/main/' > /etc/apk/repositories && \
    echo 'https://mirrors.aliyun.com/alpine/v3.10/community/' >> /etc/apk/repositories
COPY --from=build /workspace/app /usr/local/bin/
WORKDIR /app
RUN chmod +x /usr/local/bin/app \
    && apk update && apk add --no-cache tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \ 
    && echo "Asia/Shanghai" > /etc/timezone 
EXPOSE 80
CMD app