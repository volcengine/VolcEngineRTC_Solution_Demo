package handler

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_handler"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/login/login_service"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/login/login_handler"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/endpoint"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type EventHandlerDispatch struct {
	mws      []endpoint.Middleware
	eps      endpoint.Endpoint
	handlers map[string]endpoint.Endpoint
}

func NewEventHandlerDispatch() *EventHandlerDispatch {
	ehd := &EventHandlerDispatch{
		mws:      make([]endpoint.Middleware, 0),
		handlers: make(map[string]endpoint.Endpoint),
	}
	ehd.init()
	return ehd
}

func (ehd *EventHandlerDispatch) Handle(ctx context.Context, params *public.EventParam) (resp interface{}, err error) {
	return ehd.eps(ctx, params)
}

func (ehd *EventHandlerDispatch) init() {
	ehd.initHandlers()
	ehd.initMws()
	ehd.buildInvokeChain()
}

func (ehd *EventHandlerDispatch) initHandlers() {
	//disconnect
	ehd.register("disconnect", disconnectHandler)

	//login
	loginHandler := login_handler.NewEventHandler()
	ehd.register("passwordFreeLogin", loginHandler.PasswordFreeLogin)
	ehd.register("verifyLoginToken", loginHandler.VerifyLoginToken)
	ehd.register("changeUserName", loginHandler.ChangeUserName)
	ehd.register("joinRTM", loginHandler.JoinRtm)

	//voice_chat
	svcHandler := svc_handler.NewEventHandler()
	ehd.register("svcAgreeApply", svcHandler.AgreeApply)
	ehd.register("svcApplyInteract", svcHandler.ApplyInteract)
	ehd.register("svcFinishInteract", svcHandler.FinishInteract)
	ehd.register("svcFinishLive", svcHandler.FinishLive)
	ehd.register("svcGetActiveLiveRoomList", svcHandler.GetActiveLiveRoomList)
	ehd.register("svcGetApplyAudienceList", svcHandler.GetApplyAudienceList)
	ehd.register("svcGetAudienceList", svcHandler.GetAudienceList)
	ehd.register("svcInviteInteract", svcHandler.InviteInteract)
	ehd.register("svcJoinLiveRoom", svcHandler.JoinLiveRoom)
	ehd.register("svcLeaveLiveRoom", svcHandler.LeaveLiveRoom)
	ehd.register("svcManageInteractApply", svcHandler.ManageInteractApply)
	ehd.register("svcManageSeat", svcHandler.ManageSeat)
	ehd.register("svcReplyInvite", svcHandler.ReplyInvite)
	ehd.register("svcSendMessage", svcHandler.SendMessage)
	ehd.register("svcStartLive", svcHandler.StartLive)
	ehd.register("svcUpdateMediaStatus", svcHandler.UpdateMediaStatus)
	ehd.register("svcReconnect", svcHandler.Reconnect)
	ehd.register("svcClearUser", svcHandler.ClearUser)
	//other
}

func (ehd *EventHandlerDispatch) register(eventName string, handlerFunc endpoint.Endpoint) {
	ehd.handlers[eventName] = handlerFunc
}

func (ehd *EventHandlerDispatch) initMws() {
	ehd.mws = append(ehd.mws, mwCheckParam)
	ehd.mws = append(ehd.mws, mwCheckLogin)
}

func (ehd *EventHandlerDispatch) buildInvokeChain() {
	handler := ehd.invokeHandleEndpoint()
	ehd.eps = endpoint.Chain(ehd.mws...)(handler)
}

func (ehd *EventHandlerDispatch) invokeHandleEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
		f, ok := ehd.handlers[param.EventName]
		if !ok {
			return nil, custom_error.InternalError(errors.New("event not exist"))
		}
		return f(ctx, param)
	}
}

func mwCheckParam(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
		sourceApi, _ := ctx.Value(public.CtxSourceApi).(string)
		switch sourceApi {
		case "http":
			if param.EventName == "" || param.Content == "" || param.DeviceID == "" {
				return nil, custom_error.ErrInput
			}
		case "rtm":
			if param.AppID == "" || param.UserID == "" || param.EventName == "" || param.Content == "" || param.DeviceID == "" {
				return nil, custom_error.ErrInput
			}
		case "rpc":

		default:
			return nil, custom_error.InternalError(errors.New("source api missing"))
		}
		return next(ctx, param)
	}

}

type checkParam struct {
	LoginToken string `json:"login_token"`
}

func mwCheckLogin(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
		userService := login_service.GetUserService()
		sourceApi, _ := ctx.Value(public.CtxSourceApi).(string)
		switch sourceApi {
		case "rtm":
			p := &checkParam{}
			err = json.Unmarshal([]byte(param.Content), p)
			if err != nil {
				logs.CtxError(ctx, "json unmarshal failed,error:%s", err)
				return nil, custom_error.InternalError(err)
			}
			err = userService.CheckLoginToken(ctx, p.LoginToken)
			if err != nil {
				logs.CtxError(ctx, "login_token expired")
				return nil, err
			}
			loginUserID := userService.GetUserID(ctx, p.LoginToken)
			if loginUserID != param.UserID {
				return nil, custom_error.ErrorTokenUserNotMatch
			}
		}

		return next(ctx, param)
	}
}
