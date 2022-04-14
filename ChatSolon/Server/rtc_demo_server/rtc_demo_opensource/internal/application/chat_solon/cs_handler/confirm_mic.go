package cs_handler

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type confirmMicReq struct {
	LoginToken string `json:"login_token"`
	RoomID     string `json:"room_id"`
	UserID     string `json:"user_id"`
}

type confirmMicResp struct {
}

func (eh *EventHandler) ConfirmMic(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	var p confirmMicReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	if p.RoomID == "" || p.UserID == "" {
		logs.CtxWarn(ctx, "input format error, params: %v", p)
		return nil, custom_error.ErrInput
	}

	interactService, err := cs_service.NewInteractServiceByRoomID(ctx, p.RoomID)
	if err != nil {
		logs.CtxError(ctx, "get interact cs_service failed,error:%s", err)
		return nil, err
	}
	if interactService == nil {
		logs.CtxError(ctx, "get interact cs_service failed,error:%s", err)
		return nil, custom_error.InternalError(errors.New("interact is not exist"))
	}
	err = interactService.Confirm(ctx, p.UserID)
	if err != nil {
		logs.CtxError(ctx, "invite failed,error:%s", err)
		return nil, err
	}

	return nil, nil

}
