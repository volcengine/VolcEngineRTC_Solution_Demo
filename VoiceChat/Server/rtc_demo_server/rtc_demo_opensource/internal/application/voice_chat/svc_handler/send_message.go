package svc_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type sendMessageReq struct {
	RoomID     string `json:"room_id"`
	UserID     string `json:"user_id"`
	Message    string `json:"message"`
	LoginToken string `json:"login_token"`
}

type sendMessageResp struct {
}

func (eh *EventHandler) SendMessage(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	logs.CtxInfo(ctx, "liveCreateLive param:%+v", param)
	var p sendMessageReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	//todo 敏感词检测

	userFactory := svc_service.GetUserFactory()

	user, err := userFactory.GetActiveUserByRoomIDUserID(ctx, p.RoomID, p.UserID)
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return nil, err
	}
	if user == nil {
		logs.CtxError(ctx, "user is not exist,error:%s", err)
		return nil, custom_error.ErrUserNotExist
	}

	informer := svc_service.GetInformService()
	data := &svc_service.InformMessage{
		UserInfo: user,
		Message:  p.Message,
	}
	informer.BroadcastRoom(ctx, p.RoomID, svc_service.OnMessage, data)

	return nil, nil
}
