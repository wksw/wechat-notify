#!/bin/bash

ROOTDIR=$(cd $(dirname $0);pwd)
source $ROOTDIR/deploy_config.sh

# load wechat
docker load -i $ROOTDIR/wechat.tar.gz

VERSION=$(cat version)


if [ "$(docker ps -a |grep wechat_server)" = "" ];then
	docker run -d \
		--name wechat_server \
		--net wksw \
		--ip 172.172.0.7 \
		-p 8080:8080 \
		-v $ROOTDIR/app.conf:/opt/wechat/conf/app.conf:ro \
		--restart always \
		--log-opt 'max-size=100m' \
		dockerhub.wksw.com/wksw/wechat:$VERSION
else 
	docker run -d \
		--net wksw \
		-v $ROOTDIR/app.conf:/opt/wechat/conf/app.conf:ro \
		--restart always \
		--log-opt 'max-size=100m' \
		dockerhub.wksw.com/wksw/wechat:$VERSION
fi
