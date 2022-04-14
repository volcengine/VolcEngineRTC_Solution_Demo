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

type joinMeetingReq struct {
	LoginToken string `json:"login_token"`
	UserName   string `json:"user_name"`
	UserID     string `json:"user_id"`
	RoomID     string `json:"room_id"`
}

type joinMeetingResp struct {
	Token string                `json:"token"`
	Info  *cs_models.RoomInfo   `json:"info"`
	Users []*cs_models.UserInfo `json:"users"`
}

func (eh *EventHandler) JoinMeeting(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	var p joinMeetingReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	if p.RoomID == "" || p.UserName == "" || p.UserID == "" {
		logs.CtxWarn(ctx, "input format error, params: %v", p)
		return nil, custom_error.ErrInput
	}

	roomService, err := cs_service.NewRoomServiceByRoomID(ctx, p.RoomID)
	if err != nil {
		return nil, err
	}

	err = roomService.Join(ctx, p.UserID, p.UserName)
	if err != nil {
		return nil, err
	}

	userService, err := cs_service.NewUserServiceByRoomIDUserID(ctx, p.RoomID, p.UserID)
	if err != nil {
		logs.CtxError(ctx, "get user service failed,error:%s", err)
		return nil, err
	}

	informer := cs_service.GetInformService()
	informer.BroadcastRoom(ctx, p.RoomID, cs_models.OnCsJoinMeeting, cs_service.User2CsUserInfo(userService.GetUser()))

	tokenParam := &token.GenerateParam{
		AppID:        config.Configs().CsAppID,
		AppKey:       config.Configs().CsAppKey,
		RoomID:       p.RoomID,
		UserID:       p.UserID,
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

	resp = &joinMeetingResp{
		Token: rtcToken,
		Info:  cs_service.Room2CsRoomInfo(ctx, roomService.GetRoom()),
		Users: csUsers,
	}
	return resp, nil

}
