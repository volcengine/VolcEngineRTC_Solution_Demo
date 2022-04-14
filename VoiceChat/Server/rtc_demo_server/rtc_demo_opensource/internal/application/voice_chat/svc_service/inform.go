package svc_service

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/response"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/rtc_openapi"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type InformEvent string

const (
	OnAudienceJoinRoom  InformEvent = "svcOnAudienceJoinRoom"
	OnAudienceLeaveRoom InformEvent = "svcOnAudienceLeaveRoom"
	OnFinishLive        InformEvent = "svcOnFinishLive"
	OnInviteInteract    InformEvent = "svcOnInviteInteract"
	OnApplyInteract     InformEvent = "svcOnApplyInteract"
	OnInviteResult      InformEvent = "svcOnInviteResult"
	OnJoinInteract      InformEvent = "svcOnJoinInteract"
	OnFinishInteract    InformEvent = "svcOnFinishInteract"
	OnMessage           InformEvent = "svcOnMessage"
	OnMediaStatusChange InformEvent = "svcOnMediaStatusChange"
	OnMediaOperate      InformEvent = "svcOnMediaOperate"
	OnSeatStatusChange  InformEvent = "svcOnSeatStatusChange"
	OnClearUser         InformEvent = "svcOnClearUser"
)

type InformGeneral struct {
	RoomID   string `json:"room_id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
}

type InformJoinRoom struct {
	UserInfo      *User `json:"user_info"`
	AudienceCount int   `json:"audience_count"`
}

type InformLeaveRoom struct {
	UserInfo      *User `json:"user_info"`
	AudienceCount int   `json:"audience_count"`
}

type InformFinishLive struct {
	RoomID string `json:"room_id"`
	Type   int    `json:"type"`
}

type InformInviteInteract struct {
	HostInfo *User `json:"host_info"`
	SeatID   int   `json:"seat_id"`
}

type InformApplyInteract struct {
	UserInfo *User `json:"user_info"`
	SeatID   int   `json:"seat_id"`
}

type InformInviteResult struct {
	UserInfo *User `json:"user_info"`
	Reply    int   `json:"reply"`
}

type InformJoinInteract struct {
	UserInfo *User `json:"user_info"`
	SeatID   int   `json:"seat_id"`
}

type InformFinishInteract struct {
	UserInfo *User `json:"user_info"`
	SeatID   int   `json:"seat_id"`
	Type     int   `json:"type"`
}

type InformMessage struct {
	UserInfo *User  `json:"user_info"`
	Message  string `json:"message"`
}

type InformUpdateMediaStatus struct {
	UserInfo *User `json:"user_info"`
	SeatID   int   `json:"seat_id"`
	Mic      int   `json:"mic"`
}

type InformMediaOperate struct {
	Mic int `json:"mic"`
}

type InformSeatStatusChange struct {
	SeatID int `json:"seat_id"`
	Type   int `json:"type"`
}

type InformService struct {
	AppID    string
	userRepo UserRepo
}

var informService *InformService

func GetInformService() *InformService {
	if informService == nil {
		informService = &InformService{
			AppID:    config.Configs().SvcAppID,
			userRepo: GetUserRepo(),
		}
	}
	return informService
}

func (is *InformService) BroadcastRoom(ctx context.Context, roomID string, event InformEvent, data interface{}) {
	instance := rtc_openapi.GetInstance()
	err := instance.RtmSendBroadcast(ctx, is.AppID, roomID, response.NewInformToClient(string(event), data))
	if err != nil {
		logs.CtxError(ctx, "rtm send broad cast failed,event:%s,data:%#v,error:%s", event, data, err)
	}
}

func (is *InformService) UnicastUser(ctx context.Context, roomID, userID string, event InformEvent, data interface{}) {
	instance := rtc_openapi.GetInstance()
	err := instance.RtmSendRoomUnicast(ctx, is.AppID, roomID, userID, response.NewInformToClient(string(event), data))
	if err != nil {
		logs.CtxError(ctx, "rtm send broad cast failed,event:%s,data:%#v,error:%s", event, data, err)
	}
}
