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

type teacherGetGroupStudentsReq struct {
	LoginToken  string `json:"login_token"`
	UserID      string `json:"user_id"`
	GroupRoomID string `json:"group_room_id"`
}

func teacherGetGroupStudents(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("teacherGetStudentOnMicList:%+v", param)
	var p teacherGetGroupStudentsReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//校验参数
	if p.GroupRoomID == "" || p.UserID == "" {
		logs.Warnf("input group_room_id or user_id error, params: %v", p)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	parentRoomID := edu_models.GetParentRoomID(ctx, p.GroupRoomID)
	room, err := edu_models.GetRoom(ctx, parentRoomID)
	if err != nil {
		logs.Errorf("get room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrRoomNotExist)
		return
	}
	users, err := room.ListGroupRoomStudents(ctx, p.GroupRoomID)
	if err != nil {
		logs.Errorf("get group users failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	service_utils.Push2Client(ctx, param, err, users)
}
