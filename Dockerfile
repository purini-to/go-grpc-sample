FROM alpine

ENV TZ=Asia/Tokyo
ENV LANG=ja_JP.UTF-8
ENV LANGUAGE=ja_JP.UTF-8
ENV LC_ALL=ja_JP.UTF-8

ENV APP_DIR /var/lib/grpc-sample

WORKDIR $APP_DIR

ADD grpc-sample ./bin/

EXPOSE 8080 6565

ENTRYPOINT ["/var/lib/grpc-sample/bin/grpc-sample"]
