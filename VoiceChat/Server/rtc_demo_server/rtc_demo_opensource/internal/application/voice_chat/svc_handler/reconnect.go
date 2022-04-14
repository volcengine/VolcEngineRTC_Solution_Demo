package svc_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/login/login_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_service"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type reconnectReq struct {
	LoginToken string `json:"login_token"`
}

type reconnectResp struct {
	RoomInfo *svc_service.Room             `json:"room_info"`
	UserInfo *svc_service.User             `json:"user_info"`
	HostInfo *svc_service.User             `json:"host_info"`
	SeatList map[int]*svc_service.SeatInfo `json:"seat_list"`
	RtcToken string                        `json:"rtc_token"`
}

func (eh *EventHandler) Reconnect(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	logs.CtxInfo(ctx, "svcReconect param:%+v", param)
	var p reconnectReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	loginUserService := login_service.GetUserService()
	userID := loginUserService.GetUserID(ctx, p.LoginToken)

	roomFactory := svc_service.GetRoomFactory()
	userFactory := svc_service.GetUserFactory()
	seatFactory := svc_service.GetSeatFactory()

	user, err := userFactory.GetActiveUserByUserID(ctx, userID)
	if err != nil || user == nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
		return nil, custom_error.ErrUserIsInactive
	}
	user.Reconnect(param.DeviceID)
	err = userFactory.Save(ctx, user)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s")
		return nil, err
	}

	room, err := roomFactory.GetRoomByRoomID(ctx, user.GetRoomID())
	if err != nil {
		logs.CtxError(ctx, "get room failed,error:%s", err)
		return nil, err
	}
	if room == nil {
		logs.CtxError(ctx, "room is not exist")
		return nil, custom_error.ErrRoomNotExist
	}
	host, err := userFactory.GetActiveUserByRoomIDUserID(ctx, room.GetRoomID(), room.GetHostUserID())
	if err != nil {
		logs.CtxError(ctx, "get user failed,error:%s", err)
	}

	seats, err := seatFactory.GetSeatsByRoomID(ctx, user.GetRoomID())
	if err != nil {
		logs.CtxError(ctx, "get seats failed,error:%s", err)
	}
	seatList := make(map[int]*svc_service.SeatInfo)
	for _, seat := range seats {
		if seat.GetOwnerUserID() != "" {
			u, _ := userFactory.GetActiveUserByRoomIDUserID(ctx, seat.GetRoomID(), seat.GetOwnerUserID())
			seatList[seat.SeatID] = &svc_service.SeatInfo{
				Status:    seat.Status,
				GuestInfo: u,
			}
		} else {
			seatList[seat.SeatID] = &svc_service.SeatInfo{
				Status:    seat.Status,
				GuestInfo: nil,
			}
		}
	}

	resp = &reconnectResp{
		RoomInfo: room,
		HostInfo: host,
		UserInfo: user,
		SeatList: seatList,
		RtcToken: room.GenerateToken(ctx, user.GetUserID()),
	}

	return resp, nil

}
