FROM ubuntu:18.04

LABEL maintainer="1024769485@qq.com"

WORKDIR /server/
COPY publisher /server/

CMD ["/server/publisher", "-alsologtostderr=true", "-v", "4"]
