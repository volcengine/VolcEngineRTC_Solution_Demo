package handler

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_handler"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
)

func disconnectHandler(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	svcHandler := svc_handler.NewEventHandler()
	svcHandler.Disconnect(ctx, param)

	return nil, nil
}
