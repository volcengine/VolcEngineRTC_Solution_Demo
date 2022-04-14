package svc_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type getActiveLiveRoomListReq struct {
	LoginToken string `json:"login_token"`
}

type getActiveLiveRoomListResp struct {
	RoomList []*svc_service.Room `json:"room_list"`
}

func (eh *EventHandler) GetActiveLiveRoomList(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	logs.CtxInfo(ctx, "svcGetActiveLiveRoomList param:%+v", param)
	var p getActiveLiveRoomListReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	roomFactory := svc_service.GetRoomFactory()
	roomList, err := roomFactory.GetActiveRoomList(ctx, true)
	if err != nil {
		logs.CtxError(ctx, "get active room list failed,error:%s", err)
		return nil, err
	}

	resp = &getActiveLiveRoomListResp{
		RoomList: roomList,
	}

	return resp, nil
}
