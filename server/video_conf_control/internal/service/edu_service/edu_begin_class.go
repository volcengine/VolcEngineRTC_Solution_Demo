package edu_service

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type beginOrEndClassReq struct {
	RoomID     string `json:"room_id"`
	UserID     string `json:"user_id"`
	LoginToken string `json:"login_token"`
}

func beginClass(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("beginClass:%+v", param)
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
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrRoomNotExist)
		return
	}
	if room.GetTeacherUserID() != p.UserID {
		logs.Errorf("user is not teacher of this class")
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrUserRoleNotMatch)
		return
	}

	//开始上课
	room.BeginClass(ctx)
	err = room.Save(ctx)
	if err != nil {
		logs.Errorf("begin class failed,room:%v,error:%s", room, err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	edu_models.StartRecord(ctx, config.Config.AppID, room.GetRoomID(), room.GetTeacherUserID(), room.GetRoomName())
	room.InformRoom(ctx, edu_models.OnBeginClass, &edu_models.NoticeRoom{RoomID: room.GetRoomID()})
	service_utils.Push2Client(ctx, param, err, &edu_models.NoticeRoom{RoomID: room.GetRoomID()})
}
