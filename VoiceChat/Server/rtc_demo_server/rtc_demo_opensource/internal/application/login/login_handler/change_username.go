package login_handler

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/login/login_service"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type userNameReq struct {
	LoginToken string `json:"login_token"`
	UserName   string `json:"user_name"`
}

func (h *EventHandler) ChangeUserName(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	var p userNameReq
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

	if p.UserName == "" {
		logs.CtxWarn(ctx, "input format error, params: %v", p)
		return nil, custom_error.ErrInput
	}

	userID := userService.GetUserID(ctx, p.LoginToken)
	if userID == "" {
		logs.CtxError(ctx, "get empty userID")
		err = errors.New("empty userID")
		return nil, custom_error.InternalError(err)
	}

	err = userService.SetUserName(ctx, userID, p.UserName)
	if err != nil {
		logs.CtxError(ctx, "failed to set username, err: %v", err)
		return nil, err
	}

	return nil, nil
}
