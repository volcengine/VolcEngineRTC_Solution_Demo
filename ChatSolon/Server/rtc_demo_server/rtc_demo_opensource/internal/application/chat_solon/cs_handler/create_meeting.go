package cs_handler

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_service"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/token"
)

const (
	maxUserNameLen = 64
	maxRoomNameLen = 128
	retryDelay     = 8 * time.Millisecond
	maxRetryDelay  = 128 * time.Millisecond
	maxRetryNum    = 10
)

type createMeetingReq struct {
	LoginToken string `json:"login_token"`
	UserID     string `json:"user_id"`
	UserName   string `json:"user_name"`
	RoomName   string `json:"room_name"`
}

type createMeetingResp struct {
	Token string                `json:"token"`
	Info  *cs_models.RoomInfo   `json:"info"`
	Users []*cs_models.UserInfo `json:"users"`
}

func (eh *EventHandler) CreateMeeting(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	var p createMeetingReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	if p.UserID == "" || p.UserName == "" || p.RoomName == "" {
		logs.CtxWarn(ctx, "input format error, params: %v", p)
		return nil, custom_error.ErrInput
	}

	roomName, err := url.QueryUnescape(p.RoomName)
	if err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	if len(p.UserName) > maxUserNameLen || len(roomName) > maxRoomNameLen {
		logs.CtxWarn(ctx, "input format error, params: %v", p)
		return nil, custom_error.ErrInput
	}

	hall := cs_service.GetHall()
	roomID, err := hall.GenerateRoomIDWithRetry(ctx)
	if err != nil {
		logs.CtxError(ctx, "generate room id failed,error:%s", err)
		return nil, err
	}
	userID := param.DeviceID

	err = hall.CleanUserResidual(ctx, userID)
	if err != nil {
		logs.CtxError(ctx, "clean user residual failed,error:%s", err)
		return nil, err
	}

	createRoomParam := &cs_service.CreateRoomParam{
		AppID:        config.Configs().CsAppID,
		RoomID:       roomID,
		RoomName:     p.RoomName,
		HostUserID:   p.UserID,
		HostUserName: p.UserName,
	}
	room, host, err := hall.CreateRoom(ctx, createRoomParam)
	if err != nil {
		logs.CtxError(ctx, "create room failed,error:%s", err)
		return nil, err
	}

	roomService := cs_service.NewRoomService(ctx, room)
	_, err = roomService.CreateInteract(ctx, roomID)
	if err != nil {
		return nil, err
	}
	err = roomService.Start(ctx)
	if err != nil {
		logs.CtxError(ctx, "start room cs_service failed,error:%s", err)
		return nil, err
	}

	tokenParam := &token.GenerateParam{
		AppID:        config.Configs().CsAppID,
		AppKey:       config.Configs().CsAppKey,
		RoomID:       room.GetRoomID(),
		UserID:       p.UserID,
		ExpireAt:     7 * 24 * 3600,
		CanPublish:   true,
		CanSubscribe: true,
	}
	rtcToken, err := token.GenerateToken(tokenParam)
	resp = &createMeetingResp{
		Token: rtcToken,
		Info:  cs_service.Room2CsRoomInfo(ctx, room),
		Users: []*cs_models.UserInfo{cs_service.User2CsUserInfo(host)},
	}

	return resp, nil
}
