#!/bin/bash

ROOTDIR=$(cd $(dirname $0);pwd)
source $ROOTDIR/deploy_config.sh
docker exec mariadb_server mysql -uroot -p$MYSQL_ROOT_PASSWORD -e "create database if not exists keystone;grant all privileges on keystone.* to 'keystone'@'%' identified by '$MYSQL_KEYSTONE_PASSWORD';flush privileges;"
docker exec mariadb_server mysql -uroot -p$MYSQL_ROOT_PASSWORD -e "create database if not exists notify;grant all privileges on notify.* to 'notify'@'%' identified by '$MYSQL_NOTIFY_PASSWORD';flush privileges"