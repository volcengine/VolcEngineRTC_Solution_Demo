package http

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/pkg/record"
)

func Run() {
	router := fasthttprouter.New()

	router.POST("/RecordCallBack", record.HandleRecordCallbackHttp)
	fasthttp.ListenAndServe(":12345", router.Handler)
}
