package cs_handler

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/token"
)

type reconnectMeetingReq struct {
	LoginToken string `json:"login_token"`
}

type reconnectMeetingResp struct {
	Token string                `json:"token"`
	Info  *cs_models.RoomInfo   `json:"info"`
	Users []*cs_models.UserInfo `json:"users"`
}

func (eh *EventHandler) Reconnect(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	var p reconnectMeetingReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	if param.RoomID == "" || param.UserID == "" {
		logs.CtxWarn(ctx, "input format error, params: %v", p)
		return nil, custom_error.ErrInput
	}

	roomService, err := cs_service.NewRoomServiceByRoomID(ctx, param.RoomID)
	if err != nil {
		return nil, err
	}

	userService, err := cs_service.NewReconnectUserServiceByRoomIDUserID(ctx, param.RoomID, param.UserID)
	if err != nil {
		logs.CtxError(ctx, "get user service failed,error:%s", err)
		return nil, err
	}

	err = userService.Reconnect(ctx)
	if err != nil {
		logs.CtxError(ctx, "reconnect failed,error:%s", err)
		return nil, err
	}

	tokenParam := &token.GenerateParam{
		AppID:        config.Configs().CsAppID,
		AppKey:       config.Configs().CsAppKey,
		RoomID:       param.RoomID,
		UserID:       param.UserID,
		ExpireAt:     7 * 24 * 3600,
		CanPublish:   true,
		CanSubscribe: true,
	}
	rtcToken, err := token.GenerateToken(tokenParam)

	users, err := roomService.ListUsers(ctx)
	if err != nil {
		logs.CtxError(ctx, "get audiences failed,error:%s", err)
		return nil, err
	}

	csUsers := make(cs_models.UserInfoSlice, 0)
	for _, u := range users {
		csUsers = append(csUsers, cs_service.User2CsUserInfo(u))
	}
	sort.Sort(csUsers)

	resp = &reconnectMeetingResp{
		Token: rtcToken,
		Info:  cs_service.Room2CsRoomInfo(ctx, roomService.GetRoom()),
		Users: csUsers,
	}
	return resp, nil
}
