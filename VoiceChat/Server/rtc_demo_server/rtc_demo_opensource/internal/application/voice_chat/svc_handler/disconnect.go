package svc_handler

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

func (eh *EventHandler) Disconnect(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	logs.CtxInfo(ctx, "svc disconnect,param:%#v", param)

	roomService := svc_service.GetRoomService()
	roomService.Disconnect(ctx, param.RoomID, param.UserID)

	return nil, nil
}
