package cs_handler

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_service"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type getRaiseHandsReq struct {
	RoomID     string `json:"room_id"`
	LoginToken string `json:"login_token"`
}

type getRaiseHandsResp struct {
	Users []*cs_models.UserInfo `json:"users"`
}

func (eh *EventHandler) GetRaiseHands(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	var p getRaiseHandsReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	if p.RoomID == "" {
		logs.CtxWarn(ctx, "input format error, params: %v", p)
		return nil, custom_error.ErrInput
	}

	roomService, err := cs_service.NewRoomServiceByRoomID(ctx, p.RoomID)
	if err != nil {
		return nil, err
	}
	users, err := roomService.ListAudiencesByInteractStatus(ctx, []int{cs_models.UserInteractStatusRaiseHands, cs_models.UserInteractStatusOnMicrophone})
	if err != nil {
		return nil, err
	}

	csUsers := make(cs_models.UserInfoRaiseHandsSlice, 0)
	for _, u := range users {
		csUsers = append(csUsers, cs_service.User2CsUserInfo(u))
	}
	sort.Sort(csUsers)

	resp = &getAudiencesResp{
		Users: csUsers,
	}

	return resp, nil

}
