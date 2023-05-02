FROM golang:latest
LABEL authors="robertwang"



WORKDIR /go/src/app
COPY . .

#原始方式：直接镜像内打包编译
RUN export GO111MODULE=auto && export GOPROXY=https://goproxy.cn && go mod tidy
RUN go build -o go_gateway

CMD ./bin/go_gateway -config=./conf/prod/ -endpoint=server