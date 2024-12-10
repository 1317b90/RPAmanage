
FROM golang:1.23-alpine

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 设置中国时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
  
WORKDIR /work
 
COPY go.* ./
RUN go mod download
 
# COPY All things inside the project and build
COPY . .
RUN go build -o /work/build/myapp .
 
EXPOSE 8080
ENTRYPOINT [ "/work/build/myapp" ]