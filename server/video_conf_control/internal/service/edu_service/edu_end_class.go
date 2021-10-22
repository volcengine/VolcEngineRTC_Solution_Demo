package edu_service

import (
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"

	"context"
	"encoding/json"
	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

func endClass(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("endClass:%+v", param)
	var p beginOrEndClassReq
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
		logs.Errorf("get room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}
	if p.UserID != room.GetTeacherUserID() {
		logs.Errorf("user is not teacher of this class")
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrUserRoleNotMatch)
		return
	}

	room.InformRoom(ctx, edu_models.OnEndClass, &edu_models.NoticeRoom{RoomID: room.GetRoomID()})
	edu_models.StopRecord(ctx, config.Config.AppID, room.GetRoomID(), room.GetTeacherUserID())

	err = room.EndClass(ctx)
	if err != nil {
		logs.Errorf("end class failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}
	err = room.Save(ctx)
	if err != nil {
		logs.Errorf("save failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	//返回
	service_utils.Push2Client(ctx, param, err, &edu_models.NoticeRoom{RoomID: p.RoomID})
}
