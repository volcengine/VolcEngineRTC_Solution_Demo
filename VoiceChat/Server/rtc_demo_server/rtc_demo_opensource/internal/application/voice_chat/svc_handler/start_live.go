package svc_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type startLiveReq struct {
	UserID              string `json:"user_id"`
	UserName            string `json:"user_name"`
	RoomName            string `json:"room_name"`
	BackgroundImageName string `json:"background_image_name"`
	LoginToken          string `json:"login_token"`
}

type startLiveResp struct {
	RoomInfo *svc_service.Room `json:"room_info"`
	UserInfo *svc_service.User `json:"user_info"`
	RtcToken string            `json:"rtc_token"`
}

func (eh *EventHandler) StartLive(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	logs.CtxInfo(ctx, "svcStartLive param:%+v", param)
	var p startLiveReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	//校验参数
	if p.UserID == "" || p.UserName == "" || p.RoomName == "" {
		logs.CtxError(ctx, "input error, param:%v", p)
		return nil, custom_error.ErrInput
	}

	//todo 敏感词检测pkg

	//创建房间
	roomService := svc_service.GetRoomService()

	room, host, err := roomService.CreateRoom(ctx, p.RoomName, p.BackgroundImageName, p.UserID, p.UserName, param.DeviceID)
	if err != nil {
		logs.CtxError(ctx, "create room failed,error:%s", err)
		return nil, err
	}

	//开始直播
	err = roomService.StartLive(ctx, room.GetRoomID())
	if err != nil {
		logs.CtxError(ctx, "room start live failed,error:%s", err)
		return nil, err
	}

	resp = &startLiveResp{
		RoomInfo: room,
		UserInfo: host,
		RtcToken: room.GenerateToken(ctx, host.GetUserID()),
	}

	return resp, nil
}
