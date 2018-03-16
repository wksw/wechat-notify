#!/bin/bash

ROOTDIR=$(cd $(dirname $0);pwd)

VERSION=$(cat $ROOTDIR/../version)
cd $ROOTDIR/../ && go build && docker build -t dockerhub.wksw.com/wksw/wechat:$VERSION .
# cd $ROOTDIR/../../notify_frontend && docker build -t dockerhub.wksw.com/wksw/notify_frontend:$VERSION .

cd $ROOTDIR
/bin/bash release.sh

rm -rf $ROOTDIR/notify