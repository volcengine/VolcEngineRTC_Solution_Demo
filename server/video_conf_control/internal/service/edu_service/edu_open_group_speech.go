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

func openGroupSpeech(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("openGroupSpeech:%+v", param)
	var p groupSpeechReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//校验参数
	if p.RoomID == "" || p.UserID == "" {
		logs.Warnf("input room_id or user_id error, params: %v", p)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	room, err := edu_models.GetRoom(ctx, p.RoomID)
	if err != nil {
		logs.Errorf("end class failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrRoomNotExist)
		return
	}
	if p.UserID != room.GetTeacherUserID() {
		logs.Errorf("user is not teacher of this class")
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrUserRoleNotMatch)
		return
	}

	room.OpenGroupSpeech()
	err = room.Save(ctx)
	if err != nil {
		logs.Errorf("save room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	room.InformRoom(ctx, edu_models.OnOpenGroupSpeech, &edu_models.NoticeRoom{RoomID: p.RoomID})
	service_utils.Push2Client(ctx, param, err, &edu_models.NoticeRoom{RoomID: p.RoomID})
}
