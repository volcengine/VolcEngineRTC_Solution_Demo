FROM golang:latest

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

LABEL version = "1.1" description="rtc_demo server"

COPY rtc_demo_opensource/ /go/src/github.com/volcengine/VolcEngineRTC_Solution_Demo/
COPY setup.sh /go/src/github.com/volcengine/VolcEngineRTC_Solution_Demo/
COPY wait-for-it-master/wait-for-it.sh /go/src/github.com/volcengine/VolcEngineRTC_Solution_Demo/

WORKDIR /go/src/github.com/volcengine/VolcEngineRTC_Solution_Demo/

CMD ["sh","setup.sh"]

