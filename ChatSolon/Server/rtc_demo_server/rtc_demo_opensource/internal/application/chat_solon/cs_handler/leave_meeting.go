package cs_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type leaveMeetingReq struct {
	RoomID     string `json:"room_id"`
	UserID     string `json:"user_id"`
	LoginToken string `json:"login_token"`
}

type leaveMeetingResp struct {
}

func (eh *EventHandler) LeaveMeeting(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	var p leaveMeetingReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	if p.RoomID == "" || p.UserID == "" {
		logs.CtxWarn(ctx, "input format error, params: %v", p)
		return nil, custom_error.ErrInput
	}

	roomService, err := cs_service.NewRoomServiceByRoomID(ctx, p.RoomID)
	if err != nil {
		return nil, err
	}

	userService, err := cs_service.NewUserServiceByRoomIDUserID(ctx, p.RoomID, p.UserID)
	if err != nil {
		return nil, err
	}
	user := userService.GetUser()

	informer := cs_service.GetInformService()
	if user.IsHost() {
		nextHost, err := roomService.ChangeHost(ctx)
		if err != nil {
			return nil, err
		}
		if nextHost == nil {
			informer.BroadcastRoom(ctx, p.RoomID, cs_models.OnCsMeetingEnd, map[string]string{
				"room_id": p.RoomID,
			})
		} else {
			data := &cs_models.HostInfo{
				FormerHostID: user.GetUserID(),
				HostInfo:     cs_service.User2CsUserInfo(nextHost),
			}
			informer.BroadcastRoom(ctx, p.RoomID, cs_models.OnCsHostChange, data)
		}
	} else {
		err = roomService.Leave(ctx, p.UserID)
	}

	informer.BroadcastRoom(ctx, p.RoomID, cs_models.OnCsLeaveMeeting, cs_service.User2CsUserInfo(userService.GetUser()))

	return nil, nil
}
