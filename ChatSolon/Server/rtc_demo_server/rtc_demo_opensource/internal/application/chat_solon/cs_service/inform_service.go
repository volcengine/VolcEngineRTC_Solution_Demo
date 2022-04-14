package cs_service

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_repository/cs_facade"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/response"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/rtc_openapi"
)

type InformService struct {
	AppID    string
	userRepo cs_facade.RoomUserRepositoryInterface
}

var informService *InformService

func GetInformService() *InformService {
	if informService == nil {
		informService = &InformService{
			AppID:    config.Configs().CsAppID,
			userRepo: cs_facade.GetRoomUserRepository(),
		}
	}
	return informService
}

func (is *InformService) BroadcastRoom(ctx context.Context, roomID string, event cs_models.InformEvent, data interface{}) {
	instance := rtc_openapi.GetInstance()
	err := instance.RtmSendBroadcast(ctx, is.AppID, roomID, response.NewInformToClient(string(event), data))
	if err != nil {
		logs.CtxError(ctx, "rtm send broad cast failed,event:%s,data:%#v,error:%s", event, data, err)
	}
}

func (is *InformService) UnicastUser(ctx context.Context, roomID, userID string, event cs_models.InformEvent, data interface{}) {
	instance := rtc_openapi.GetInstance()
	err := instance.RtmSendRoomUnicast(ctx, is.AppID, roomID, userID, response.NewInformToClient(string(event), data))
	if err != nil {
		logs.CtxError(ctx, "rtm send broad cast failed,event:%s,data:%#v,error:%s", event, data, err)
	}
}
