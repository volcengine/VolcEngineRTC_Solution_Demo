#!/bin/bash


echo 'rtc demo build'
sh build.sh
echo 'complete rtc demo build'


echo 'wait for mysql and redis ready'
./wait-for-it.sh mysql_server:3306 -t 60
./wait-for-it.sh redis_server:6379 -t 60

echo 'rtc demo server starting'
cd output
sh bootstrap.sh
echo 'complete rtc demo server start'