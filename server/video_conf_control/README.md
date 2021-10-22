# 目录结构
```go
|____models //类
| |____edu_models 
| |____custom_error //错误码的定义和描述
| |____cs_models
| |____login_models 
| |____conn_models
| |____vc_models  //会控场景使用的类，主要为User和Room
|____service //对外提供的api
| |____conn_service
| |____service_utils
| |____cs_service
| |____vc_service
| |____audit_service
| |____login_service
| |____edu_service
|____dal //数据层
| |____redis 
| | |____login_redis
| | |____edu_redis
| | |____lock //redis锁
| | |____cs_redis
| | |____vc_redis
| |____db //mysql
| | |____vc_db
| | |____edu_db 
| | |____cs_db
| | |____conn_db
| | |____login_db
|____pkg //通用工具
| |____token //token工具
| |____video //从视频云平台获取视频
| |____record //录制相关工具，包装了对RTC后处理的openapi调用
| |____response //通用返回结构
| |____pinyin //拼音
| |____task //定时任务


简称说明：
login //账号相关
conn //connection websocket连接相关
cs //chat salon语音沙龙
vc //video conference 视频会议相关
edu //education 教育相关
```


# 环境依赖
|  地址 | 用途  |
| --- | --- |
| [github.com/volcengine/volc-sdk-golang](github.com/volcengine/volc-sdk-golang) | 视频云toB开源SDK，用来获取视频url，删除视频。 |
| [gorm.io/gorm](https://gorm.io/) | mysql |
| [github.com/go-redis/redis/v8](github.com/go-redis/redis/v8) | redis |
| [github.com/mozillazg/go-pinyin](github.com/mozillazg/go-pinyin) | 汉字转pinyin |
| [github.com/satori/go.uuid](github.com/satori/go.uuid) | uuid |
| [github.com/robfig/cron/v3](github.com/robfig/cron/v3) | 定时任务(定时清理房间、用户、录制) |
| [github.com/valyala/fasthttpgithub.com/buaazp/fasthttprouter](github.com/valyala/fasthttpgithub.com/buaazp/fasthttprouter) | http |
| [github.com/sirupsen/logrus](github.com/sirupsen/logrus) | 日志 |