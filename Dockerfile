FROM --platform=linux/arm golang:1.16-alpine
#添加工程进来
COPY . /go/src/cache
COPY ./redisarm/redis-server /home
COPY ./redisarm/redis-cli /home
COPY ./redisarm/redis.conf /home
#设置时间为北京时间
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && chmod +x /go/src/cache/cache \
    && chmod +x /home/redis-server
#暴露http端口
EXPOSE 8080 9000 6379
WORKDIR /go/src/cache
CMD  /home/redis-server /home/redis.conf && /go/src/cache/cache
