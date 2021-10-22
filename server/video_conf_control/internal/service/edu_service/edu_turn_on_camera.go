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

func turnOnCamera(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("turnOnCamera:%+v", param)
	var p turnCameraReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//校验参数
	if p.UserID == "" || p.RoomID == "" {
		logs.Warnf("input user_id or room_id error, params: %v", p)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	room, err := edu_models.GetRoom(ctx, p.RoomID)
	if err != nil {
		logs.Errorf("get room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrRoomNotExist)
		return
	}

	user, err := edu_models.GetUser(ctx, param.GetConnId())
	if err != nil {
		logs.Errorf("get user failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	user.TurnOnCamera()
	err = user.Save(ctx)
	if err != nil {
		logs.Errorf("turn on camera failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	room.InformRoom(ctx, edu_models.OnTeacherCameraOn, edu_models.NoticeRoom{RoomID: room.GetRoomID()})
	service_utils.Push2Client(ctx, param, err, &edu_models.NoticeRoom{RoomID: p.RoomID, UserID: p.UserID})
}
