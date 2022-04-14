package cs_handler

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

type getAppIDParam struct {
	LoginToken string `json:"login_token"`
}

func (eh *EventHandler) GetAppID(ctx context.Context, param *public.EventParam) (resp interface{}, err error) {
	var p getAppIDParam
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.CtxWarn(ctx, "input format error, err: %v", err)
		return nil, custom_error.ErrInput
	}

	return map[string]string{"app_id": config.Configs().CsAppID}, nil
}
