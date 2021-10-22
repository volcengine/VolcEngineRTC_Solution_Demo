module github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control

go 1.14

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0

require (
	github.com/apache/thrift v0.13.0
	github.com/buaazp/fasthttprouter v0.1.1
	github.com/cloudwego/kitex v0.0.3
	github.com/go-redis/redis/v8 v8.11.1
	github.com/jinzhu/configor v1.2.1
	github.com/mozillazg/go-pinyin v0.18.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/valyala/fasthttp v1.28.0
	github.com/volcengine/volc-sdk-golang v1.0.21
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.12
)
