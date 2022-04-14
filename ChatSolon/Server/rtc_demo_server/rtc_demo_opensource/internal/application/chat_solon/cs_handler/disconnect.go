package cs_handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

func (eh *EventHandler) Disconnect(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	logs.CtxInfo(ctx, "cs disconnect,param:%#v", param)

	userService, err := cs_service.NewUserServiceByRoomIDUserID(ctx, param.RoomID, param.UserID)
	if err != nil {
		if cerr, ok := err.(*custom_error.CustomError); ok {
			if cerr.Code() == custom_error.ErrUserNotExist.Code() {
				return nil, nil
			}
		}
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return nil, err
	}

	err = userService.Disconnect(ctx)
	if err != nil {
		logs.CtxError(ctx, "disconnect failed,error:%s", err)
		return nil, err
	}

	go func(ctx context.Context, param *public.EventParam) {
		time.Sleep(time.Duration(config.Configs().ReconnectTimeout))
		userService, err := cs_service.NewUserServiceByRoomIDUserID(ctx, param.RoomID, param.UserID)
		if err != nil {
			logs.CtxInfo(ctx, "get user failed,error:%s", err)
		}
		if userService.GetUser().GetNetStatus() == cs_entity.UserNetStatusReconnecting {
			req := leaveMeetingReq{
				RoomID: param.RoomID,
				UserID: param.UserID,
			}
			content, _ := json.Marshal(req)
			p := &public.EventParam{
				AppID:     param.AppID,
				RoomID:    param.RoomID,
				UserID:    param.UserID,
				EventName: "csLeaveMeeting",
				Content:   string(content),
			}
			eh.LeaveMeeting(ctx, p)
		}

	}(ctx, param)

	return nil, nil

}
