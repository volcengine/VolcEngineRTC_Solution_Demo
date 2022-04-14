package svc_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type updateMediaStatusReq struct {
	RoomID     string `json:"room_id"`
	UserID     string `json:"user_id"`
	Mic        int    `json:"mic"`
	LoginToken string `json:"login_token"`
}

type updateMediaStatusResp struct {
}

func (eh *EventHandler) UpdateMediaStatus(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	logs.CtxInfo(ctx, "liveCreateLive param:%+v", param)
	var p updateMediaStatusReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

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

	if p.Mic == 0 {
		user.MuteMic()
	} else {
		user.UnmuteMic()
	}
	err = userFactory.Save(ctx, user)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return nil, err
	}

	informer := svc_service.GetInformService()
	data := &svc_service.InformUpdateMediaStatus{
		UserInfo: user,
		SeatID:   user.GetSeatID(),
		Mic:      user.Mic,
	}
	informer.BroadcastRoom(ctx, p.RoomID, svc_service.OnMediaStatusChange, data)

	return nil, nil

}
