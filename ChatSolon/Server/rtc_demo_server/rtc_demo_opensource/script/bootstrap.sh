#! /usr/bin/env bash
CURDIR=$(cd $(dirname $0); pwd)
if [ "X$1" != "X" ]; then
    RUNTIME_ROOT=$1
else
    RUNTIME_ROOT=${CURDIR}
fi

PORT=$2
CONFIG_FILE=$CURDIR/conf/config.yaml

RUNTIME_LOG_ROOT=$RUNTIME_ROOT/log
export RUNTIME_LOG_ROOT=$RUNTIME_LOG_ROOT

#  service log path : $RUNTIME_LOG_ROOT/app/${svc_name}.log
if [ ! -d $RUNTIME_LOG_ROOT/app ]; then
    mkdir -p $RUNTIME_LOG_ROOT/app
fi

exec ${CURDIR}/bin/volcengine.VolcEngineRTC_Solution_Demo.rtc_demo_opensource -config="$CONFIG_FILE" -port="${PORT}"