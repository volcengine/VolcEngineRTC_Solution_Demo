y#!/usr/bin/env bash

go env -w GOPROXY=https://goproxy.cn,direct
go mod init
go mod tidy

sh build.sh
cd output
sh bootstrap.sh