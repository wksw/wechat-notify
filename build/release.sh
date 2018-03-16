#!/bin/bash

ROOTDIR=$(cd $(dirname $0);pwd)
NOW=$(date +%Y%m%d)
RELEASE_DIR=$ROOTDIR/../releases/$NOW
VERSION=$(cat $ROOTDIR/../version)


rm -rf $RELEASE_DIR ${RELEASE_DIR}.tar.gz
mkdir -p $RELEASE_DIR


# package redis
docker save -o $RELEASE_DIR/redis.tar.gz dockerhub.wksw.com/database/redis:latest

# package ssdb
docker save -o $RELEASE_DIR/ssdb.tar.gz dockerhub.wksw.com/database/ssdb:1.9.4

# package notify
docker save -o $RELEASE_DIR/wechat.tar.gz dockerhub.wksw.com/wksw/wechat:$VERSION

#package mariadb
docker save -o $RELEASE_DIR/mariadb.tar.gz dockerhub.wksw.com/database/mariadb:latest

# package keystone
docker save -o $RELEASE_DIR/keystone.tar.gz dockerhub.wksw.com/wksw/keystone:v1

# package frontend
docker save -o $RELEASE_DIR/notify_frontend.tar.gz dockerhub.wksw.com/wksw/notify_frontend:$VERSION

cp -r $ROOTDIR/../deploy/* $RELEASE_DIR
cp $ROOTDIR/../conf/app.conf $RELEASE_DIR
cp $ROOTDIR/../version $RELEASE_DIR
cp $ROOTDIR/../../notify_frontend/src/common/config/index.js $RELEASE_DIR


cd $RELEASE_DIR/../ && rm -rf ${NOW}.tar.gz && tar zcf ${NOW}.tar.gz $NOW
rm -rf $RELEASE_DIR