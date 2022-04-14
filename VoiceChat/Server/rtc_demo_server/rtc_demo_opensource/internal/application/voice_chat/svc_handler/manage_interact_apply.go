package svc_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type manageInteractApplyReq struct {
	RoomID     string `json:"room_id"`
	UserID     string `json:"user_id"`
	Type       int    `json:"type"`
	LoginToken string `json:"login_token"`
}

type manageInteractApplyResp struct {
}

func (eh *EventHandler) ManageInteractApply(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	logs.CtxInfo(ctx, "liveCreateLive param:%+v", param)
	var p manageInteractApplyReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	//校验参数
	if p.RoomID == "" || p.UserID == "" {
		logs.CtxError(ctx, "input error, param:%v", p)
		return nil, custom_error.ErrInput
	}

	roomFactory := svc_service.GetRoomFactory()
	room, err := roomFactory.GetRoomByRoomID(ctx, p.RoomID)
	if err != nil {
		logs.CtxError(ctx, "get room failed,error:%s", err)
		return nil, err
	}
	if room == nil {
		logs.CtxError(ctx, "room is not exist ")
		return nil, custom_error.ErrRoomNotExist
	}
	if room.GetHostUserID() != p.UserID {
		logs.CtxError(ctx, "user is not host of room")
		return nil, custom_error.ErrUserIsNotOwner
	}

	if p.Type == 1 {
		room.OpenApply()
	} else {
		room.CloseApply()
	}
	err = roomFactory.Save(ctx, room)
	if err != nil {
		logs.CtxError(ctx, "get room factory failed,error:%s", err)
		return nil, err
	}

	return nil, nil
}
