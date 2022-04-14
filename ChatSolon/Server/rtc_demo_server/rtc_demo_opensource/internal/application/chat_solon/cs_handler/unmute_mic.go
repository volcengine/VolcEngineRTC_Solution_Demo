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

type unmuteMicReq struct {
	LoginToken string `json:"login_token"`
	RoomID     string `json:"room_id"`
	UserID     string `json:"user_id"`
}

type unmuteMicResp struct {
}

func (eh *EventHandler) UnmuteMic(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	var p unmuteMicReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	if p.UserID == "" {
		logs.CtxWarn(ctx, "input format error, params: %v", p)
		return nil, custom_error.ErrInput
	}

	userService, err := cs_service.NewUserServiceByRoomIDUserID(ctx, p.RoomID, p.UserID)
	if err != nil {
		return nil, err
	}
	err = userService.Unmute(ctx)
	if err != nil {
		return nil, err
	}

	informer := cs_service.GetInformService()
	informer.BroadcastRoom(ctx, p.RoomID, cs_models.OnCsUnmuteMic, cs_service.User2CsUserInfo(userService.GetUser()))

	return nil, nil

}
