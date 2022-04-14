package handler

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_handler"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
)

func disconnectHandler(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	csHandler := cs_handler.NewEventHandler()
	csHandler.Disconnect(ctx, param)

	return nil, nil
}
