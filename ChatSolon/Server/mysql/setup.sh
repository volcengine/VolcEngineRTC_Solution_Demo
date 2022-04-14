#!/bin/bash
set -e

#查看mysql服务的状态，方便调试，这条语句可以删除
echo `service mysql status`

echo '1.启动mysql....'
#启动mysql
service mysql start
sleep 1
echo `service mysql status`

echo '2.创建数据库及数据表....'
#创建数据库及数据表
mysql < /mysql/rtc_demo.sql
sleep 1
echo '创建数据库及数据表....'


echo `service mysql status`
echo `mysql容器启动完毕`

tail -f /dev/null