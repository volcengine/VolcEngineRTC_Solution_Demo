package svc_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type finishLiveReq struct {
	UserID     string `json:"user_id"`
	RoomID     string `json:"room_id"`
	LoginToken string `json:"login_token"`
}

type finishLiveResp struct {
	RoomInfo *svc_service.Room `json:"room_info"`
	UserInfo *svc_service.User `json:"user_info"`
	RtcToken string            `json:"rtc_token"`
}

func (eh *EventHandler) FinishLive(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	logs.CtxInfo(ctx, "svcFinishLive param:%+v", param)
	var p finishLiveReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	//校验参数
	if p.UserID == "" || p.RoomID == "" {
		logs.CtxError(ctx, "input error, param:%v", p)
		return nil, custom_error.ErrInput
	}

	//是否是主播校验
	roomFactory := svc_service.GetRoomFactory()
	room, err := roomFactory.GetRoomByRoomID(ctx, p.RoomID)
	if err != nil || room == nil {
		logs.CtxError(ctx, "get room failed,error:%s", err)
		return nil, custom_error.ErrRoomNotExist
	}
	if room.GetHostUserID() != p.UserID {
		logs.CtxError(ctx, "user is not host of room")
		return nil, custom_error.ErrUserIsNotOwner
	}

	roomService := svc_service.GetRoomService()

	err = roomService.FinishLive(ctx, p.RoomID, svc_service.FinishTypeNormal)
	if err != nil {
		logs.CtxError(ctx, "room finish live failed,error:%s", err)
		return nil, err
	}

	return nil, nil
}
