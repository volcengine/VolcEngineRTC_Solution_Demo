package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/cmd/handler"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/cmd/api"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/redis_cli"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/task"
)

var conf = flag.String("config", "conf/config.yaml", "server config file path")

func main() {
	fmt.Println("start")
	rand.Seed(time.Now().UnixNano())

	Init()

	h := handler.NewEventHandlerDispatch()
	//start http api
	r := api.NewHttpApi(h)
	err := r.Run()
	if err != nil {
		panic("http server start failed,error:" + err.Error())
	}
}

func Init() {
	//get config
	config.InitConfig(*conf)

	logs.InitLog()

	//init db and redis
	db.Open(config.Configs().MysqlDSN)
	redis_cli.NewRedis(config.Configs().RedisAddr, config.Configs().RedisPassword)

	c := task.GetCronTask()
	c.Start()
}
