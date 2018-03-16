#!/bin/bash


ROOTDIR=$(cd $(dirname $0);pwd)

source $ROOTDIR/deploy_config.sh
VERSION=$(cat version)

# load redis
docker load -i $ROOTDIR/redis.tar.gz

# load ssdb
docker load -i $ROOTDIR/ssdb.tar.gz

#load mariadb
docker load -i $ROOTDIR/mariadb.tar.gz

#load keystone
docker load -i $ROOTDIR/keystone.tar.gz


# create network

if [ "$(docker network ls|grep wksw)" = "" ];then
	docker network create --subnet=172.172.0.1/16 wksw
fi 

# start redis
if [ "$(docker ps -a |grep redis_server)" = "" ];then
	docker run -d \
		--name redis_server \
		--net wksw \
		--ip 172.172.0.2 \
		-p 16379:6379 \
		--restart always \
		--log-opt 'max-size=100m' \
		dockerhub.wksw.com/database/redis:latest
fi


# start ssdb
if [ "$(docker ps -a |grep ssdb_server)" = "" ];then
	docker run -d \
		--name ssdb_server \
		--net wksw \
		--ip 172.172.0.3 \
		-p 18888:8888 \
		-v $ROOTDIR/ssdb/ssdb.conf:/etc/ssdb.conf \
		-v $ROOTDIR/ssdb/data:/var/lib/ssdb \
		--restart always \
		--log-opt 'max-size=100m' \
		dockerhub.wksw.com/database/ssdb:1.9.4
fi

if [ "$(docker ps -a |grep mariadb_server)" = "" ];then
	docker run -d \
		--name mariadb_server \
		--net wksw \
		--ip 172.172.0.4 \
		-p 13306:3306 \
		--restart always \
		--log-opt 'max-size=100m' \
		-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
		dockerhub.wksw.com/database/mariadb:latest \
		--character-set-server=utf8 \
		--collation-server=utf8_general_ci

fi


if [ "$(docker ps -a|grep keystone_server)" = "" ];then
	docker run -d \
		--name keystone_server \
		--net wksw \
		--ip 172.172.0.5 \
		-p 5000:5000 \
		-p 35357:35357 \
		--restart always \
		--log-opt 'max-size=100m' \
		--hostname controller \
		-e ADMIN_PASSWOR=adminwksw \
		-e HOST=172.172.0.5 \
		-v $ROOTDIR/keystone:/etc/keystone \
		--link mariadb_server:mysql \
		dockerhub.wksw.com/wksw/keystone:v1
fi

bash $ROOTDIR/deploy_wechat.sh
bash $ROOTDIR/deploy_create_database.sh

