package login_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/login/login_service"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/token"
)

type joinRtmReq struct {
	LoginToken string `json:"login_token"`
	ScenesName string `json:"scenes_name"`
}

type joinRtmResp struct {
	AppID           string `json:"app_id"`
	RtmToken        string `json:"rtm_token"`
	ServerUrl       string `json:"server_url"`
	ServerSignature string `json:"server_signature"`
}

func (h *EventHandler) JoinRtm(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	var p joinRtmReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	userService := login_service.GetUserService()

	err = userService.CheckLoginToken(ctx, p.LoginToken)
	if err != nil {
		logs.CtxWarn(ctx, "login token expiry")
		return nil, err
	}
	//todo 根据场景选appID、appKey

	userID := userService.GetUserID(ctx, p.LoginToken)

	rtcToken, err := token.GenerateToken(&token.GenerateParam{
		AppID:        config.Configs().SvcAppID,
		AppKey:       config.Configs().SvcAppKey,
		RoomID:       "",
		UserID:       userID,
		ExpireAt:     7 * 24 * 3600,
		CanPublish:   true,
		CanSubscribe: true,
	})
	if err != nil {
		logs.CtxError(ctx, "generate token failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}

	resp = &joinRtmResp{
		AppID:           config.Configs().SvcAppID,
		RtmToken:        rtcToken,
		ServerUrl:       config.Configs().ServerUrl,
		ServerSignature: p.LoginToken,
	}
	return resp, nil

}
