package svc_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type applyInteractReq struct {
	RoomID     string `json:"room_id"`
	UserID     string `json:"user_id"`
	SeatID     int    `json:"seat_id"`
	LoginToken string `json:"login_token"`
}

type applyInteractResp struct {
	IsNeedApply bool `json:"is_need_apply"`
}

func (eh *EventHandler) ApplyInteract(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	logs.CtxInfo(ctx, "svcApplyInteract param:%+v", param)
	var p applyInteractReq
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
	userFactory := svc_service.GetUserFactory()

	room, err := roomFactory.GetRoomByRoomID(ctx, p.RoomID)
	if err != nil {
		logs.CtxError(ctx, "get room failed,error:%s", err)
		return nil, err
	}
	if room == nil {
		logs.CtxError(ctx, "room is not exist")
		return nil, custom_error.ErrRoomNotExist
	}

	interactService := svc_service.GetInteractService()

	err = interactService.Apply(ctx, room.GetRoomID(), room.GetHostUserID(), p.UserID, p.SeatID)
	if err != nil {
		logs.CtxError(ctx, "invite failed,error:%s", err)
		return nil, err
	}

	if !room.IsNeedApply() {
		err = interactService.HostReply(ctx, room.GetRoomID(), room.GetHostUserID(), p.UserID, svc_service.InteractReplyTypeAccept)
		if err != nil {
			logs.CtxError(ctx, "host reply failed,error:%s", err)
			return nil, err
		}
	} else {
		informer := svc_service.GetInformService()
		user, err := userFactory.GetUserByRoomIDUserID(ctx, p.RoomID, p.UserID)
		if err != nil {
			logs.CtxError(ctx, "get user failed,error:%s", err)
		}
		data := &svc_service.InformApplyInteract{
			UserInfo: user,
			SeatID:   p.SeatID,
		}
		informer.UnicastUser(ctx, room.GetRoomID(), room.GetHostUserID(), svc_service.OnApplyInteract, data)
	}

	resp = &applyInteractResp{
		IsNeedApply: room.IsNeedApply(),
	}

	return resp, nil
}
