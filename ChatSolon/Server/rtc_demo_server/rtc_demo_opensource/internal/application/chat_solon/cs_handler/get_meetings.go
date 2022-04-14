package cs_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_service"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type getMeetingsReq struct {
	LoginToken string `json:"login_token"`
}

type getMeetingsResp struct {
	Infos []*cs_models.RoomInfo `json:"infos"`
}

func (eh *EventHandler) GetMeetings(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	var p getMeetingsReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	hall := cs_service.GetHall()
	rooms, err := hall.ListRooms(ctx)
	if err != nil {

	}
	csRooms := make([]*cs_models.RoomInfo, 0)
	for _, r := range rooms {
		csRooms = append(csRooms, cs_service.Room2CsRoomInfo(ctx, r))
	}

	resp = &getMeetingsResp{
		Infos: csRooms,
	}

	return resp, nil

}
