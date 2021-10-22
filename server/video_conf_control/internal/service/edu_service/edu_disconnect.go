package edu_service

import (
	"context"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

func Disconnect(ctx context.Context, param *vc_control.TEventParam) {

	user, err := edu_models.GetUser(ctx, param.ConnId)
	if err != nil {
		logs.Errorf("get user failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	err = user.Disconnect(ctx)
	if err != nil {
		logs.Errorf("reconnect failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, nil)
		return
	}

	service_utils.Push2Client(ctx, param, err, edu_models.NoticeRoom{})

}
