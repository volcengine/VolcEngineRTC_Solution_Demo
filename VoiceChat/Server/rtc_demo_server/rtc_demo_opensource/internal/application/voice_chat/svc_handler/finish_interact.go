package svc_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type finishInteractReq struct {
	RoomID     string `json:"room_id"`
	UserID     string `json:"user_id"`
	SeatID     int    `json:"seat_id"`
	LoginToken string `json:"login_token"`
}

type finishInteractResp struct {
}

func (eh *EventHandler) FinishInteract(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	logs.CtxInfo(ctx, "svcFinishInteract param:%+v", param)
	var p finishInteractReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	//校验参数
	if p.RoomID == "" || p.UserID == "" {
		logs.CtxError(ctx, "input error, param:%v", p)
		return nil, custom_error.ErrInput
	}

	interactService := svc_service.GetInteractService()

	err = interactService.FinishInteract(ctx, p.RoomID, p.SeatID, svc_service.InteractFinishTypeSelf)
	if err != nil {
		logs.CtxError(ctx, "finish interact failed,error:%s", err)
		return nil, err
	}

	return nil, nil
}
