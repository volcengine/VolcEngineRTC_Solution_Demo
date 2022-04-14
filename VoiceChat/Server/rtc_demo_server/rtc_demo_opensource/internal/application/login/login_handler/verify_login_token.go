package login_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/login/login_service"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type verifyTokenParam struct {
	LoginToken string `json:"login_token"`
}

func (h *EventHandler) VerifyLoginToken(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	var p verifyTokenParam
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	userService := login_service.GetUserService()

	err = userService.CheckLoginToken(ctx, p.LoginToken)
	if err != nil {
		logs.CtxWarn(ctx, "login token expiry")
		return nil, custom_error.ErrInput
	}

	return nil, nil
}
