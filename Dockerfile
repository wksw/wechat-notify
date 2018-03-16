FROM dockerhub.wksw.com/wksw/wechat:base

WORKDIR /opt/wechat
EXPOSE 8080

COPY notify /opt/wechat/

COPY conf/ /opt/wechat/conf/

CMD ["./notify"]
