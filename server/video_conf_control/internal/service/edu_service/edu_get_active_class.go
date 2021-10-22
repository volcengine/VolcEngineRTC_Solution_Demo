package edu_service

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type getActiveClassReq struct {
	LoginToken string `json:"login_token"`
}

func getActiveClass(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("getActiveClass:%+v", param)
	var p getActiveClassReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//查大房间
	rooms, err := edu_models.GetActiveRooms(ctx)
	if err != nil {
		logs.Errorf("failed to QueryActiveRooms: %s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	service_utils.Push2Client(ctx, param, err, rooms)
}
