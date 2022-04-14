package svc_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_service"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type clearUserReq struct {
	UserID     string `json:"user_id"`
	LoginToken string `json:"login_token"`
}

type clearUserResp struct {
}

func (eh *EventHandler) ClearUser(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	logs.CtxInfo(ctx, "svcClearUser param:%+v", param)
	var p clearUserReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	userFactory := svc_service.GetUserFactory()
	user, err := userFactory.GetActiveUserByUserID(ctx, p.UserID)
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	roomService := svc_service.GetRoomService()
	if user.IsHost() {
		err = roomService.FinishLive(ctx, user.GetRoomID(), svc_service.FinishTypeNormal)
		if err != nil {
			logs.CtxError(ctx, "finish live failed,error:%s", err)
			return nil, err
		}
	} else if user.IsAudience() {
		err = roomService.LeaveRoom(ctx, user.GetRoomID(), user.GetUserID())
		if err != nil {
			logs.CtxError(ctx, "leave room failed,error:%s", err)
			return nil, err
		}
	}

	//强制清除
	user, err = userFactory.GetActiveUserByUserID(ctx, p.UserID)
	if user != nil {
		user.LeaveRoom()
		userFactory.Save(ctx, user)
	}

	return nil, nil
}
