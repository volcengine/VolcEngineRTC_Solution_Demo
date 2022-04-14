package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/cmd/handler"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/response"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/rtc_openapi"
)

type HttpApi struct {
	r          *gin.Engine
	dispatcher *handler.EventHandlerDispatch
	Addr       string
	Port       string
}

func NewHttpApi(dispatcher *handler.EventHandlerDispatch) *HttpApi {
	api := &HttpApi{}
	api.r = gin.Default()
	api.dispatcher = dispatcher
	api.Addr = config.Configs().Addr
	api.Port = config.Configs().Port
	return api
}

func (api *HttpApi) Run() error {
	rr := api.r.Group("/rtc_demo")
	rr.POST("/login", api.HandleHttpLoginEvent)
	rr.POST("/rtm", api.HandleRtmOpenApiEvent)
	rr.POST("/rtm_callback", api.HandleRtmCallback)
	return api.r.Run(fmt.Sprintf("%s:%s", api.Addr, api.Port))
}

func (api *HttpApi) HandleHttpLoginEvent(httpCtx *gin.Context) {
	ctx := util.EnsureID(httpCtx)
	ctx = context.WithValue(ctx, public.CtxSourceApi, "http")

	p := &public.EventParam{}
	err := httpCtx.BindJSON(p)
	if err != nil {
		logs.CtxError(ctx, "param error,err:%s", err)
		httpCtx.String(400, "param error")
		return
	}
	logs.CtxInfo(ctx, "handle http,param:%#v", p)
	resp, err := api.dispatcher.Handle(ctx, p)
	if err != nil {
		logs.CtxError(ctx, "handle error,param:%#v,err:%s", p, err)
	}

	httpCtx.String(200, response.NewCommonResponse(ctx, "", resp, err))
	return
}

func (api *HttpApi) HandleRtmOpenApiEvent(httpCtx *gin.Context) {
	defer httpCtx.String(200, "ok")
	ctx := util.EnsureID(httpCtx)
	ctx = context.WithValue(ctx, public.CtxSourceApi, "rtm")

	p := &RtmParam{}
	err := httpCtx.BindJSON(p)
	if err != nil {
		logs.CtxError(ctx, "param error,err:%s", err)
		return
	}
	pp := &public.EventParam{}
	err = json.Unmarshal([]byte(p.Message), pp)
	if err != nil {
		logs.CtxError(ctx, "param error,err:%s", err)
		rtmReturn(ctx, pp.AppID, pp.RoomID, pp.UserID, pp.RequestID, nil, err)
		return
	}
	logs.CtxInfo(ctx, "handle rtm,param:%#v", pp)

	go func(ctx context.Context, param *public.EventParam) {
		util.CheckPanic()
		resp, err := api.dispatcher.Handle(ctx, param)
		if err != nil {
			logs.CtxError(ctx, "handle error,param:%#v,err:%s", param, err)
		}
		rtmReturn(ctx, param.AppID, param.RoomID, param.UserID, param.RequestID, resp, err)
	}(ctx, pp)

}

func rtmReturn(ctx context.Context, appID, roomID, userID, requestID string, resp interface{}, err error) {
	instance := rtc_openapi.GetInstance()

	if roomID == "" {
		err = instance.RtmSendUnicast(ctx, appID, userID, response.NewCommonResponse(ctx, requestID, resp, err))
	} else {
		err = instance.RtmSendRoomUnicast(ctx, appID, roomID, userID, response.NewCommonResponse(ctx, requestID, resp, err))
	}
	if err != nil {
		logs.CtxError(ctx, "send to rtm failed,error:%s", err)
	}
}

func (api *HttpApi) HandleRtmCallback(httpCtx *gin.Context) {
	defer httpCtx.String(200, "ok")

	ctx := util.EnsureID(httpCtx)
	ctx = context.WithValue(ctx, public.CtxSourceApi, "http")

	p := &RtmCallbackParam{}
	err := httpCtx.BindJSON(p)
	if err != nil {
		logs.CtxError(ctx, "param error,err:%s", err)
		return
	}

	pp := &public.EventParam{}

	switch p.EventType {
	case "UserLeaveRoom":
		eventData := &EventDataLeaveRoom{}
		err = json.Unmarshal([]byte(p.EventData), eventData)
		if err != nil {
			logs.CtxError(ctx, "param error,err:%s", err)
			return
		}
		if eventData.Reason != LeaveRoomReasonConnectionLost {
			return
		}
		pp.AppID = p.AppId
		pp.RoomID = eventData.RoomId
		pp.UserID = eventData.UserId
		pp.EventName = "disconnect"
	}

	logs.CtxInfo(ctx, "handle rtm callback,param:%#v", pp)
	_, err = api.dispatcher.Handle(ctx, pp)
	if err != nil {
		logs.CtxError(ctx, "handle error,param:%#v,err:%s", p, err)
	}

}
