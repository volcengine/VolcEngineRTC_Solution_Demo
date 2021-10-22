package edu_service

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/conn_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type leaveClassReq struct {
	UserID     string `json:"user_id"`
	RoomID     string `json:"room_id"`
	LoginToken string `json:"login_token"`
}

func leaveClass(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("leaveClass:%+v", param)
	var p leaveClassReq
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
		//service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrRoomNotExist)
		service_utils.Push2Client(ctx, param, nil, &edu_models.NoticeRoom{RoomID: p.RoomID, UserID: p.UserID})
		return
	}

	if p.UserID == room.GetTeacherUserID() {
		logs.Errorf("user is teacher of this class,forbid join class")
		//service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrUserRoleNotMatch)
		service_utils.Push2Client(ctx, param, nil, &edu_models.NoticeRoom{RoomID: p.RoomID, UserID: p.UserID})
		return
	}

	conn, err := conn_models.GetConnection(ctx, param.ConnId)
	if err != nil {
		logs.Errorf("failed to get connection, err: %v", err)
		//service_utils.Push2ClientWithoutReturn(ctx, param, err)
		service_utils.Push2Client(ctx, param, nil, &edu_models.NoticeRoom{RoomID: p.RoomID, UserID: p.UserID})
		return
	}

	user, err := edu_models.GetUser(ctx, conn.GetConnID())
	if err != nil || user == nil {
		logs.Errorf("get user failed,error:%s", err)
		//service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		service_utils.Push2Client(ctx, param, nil, &edu_models.NoticeRoom{RoomID: p.RoomID, UserID: p.UserID})
		return
	}

	isInteract := user.IsInteract()
	user.LeaveRoom(ctx)
	err = user.Save(ctx)
	if err != nil {
		logs.Errorf("save user failed,error:%s", err)
		//service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		service_utils.Push2Client(ctx, param, nil, &edu_models.NoticeRoom{RoomID: p.RoomID, UserID: p.UserID})
		return
	}

	if isInteract {
		room.InformRoom(ctx, edu_models.OnFinishInteract, edu_models.NoticeRoom{RoomID: room.GetRoomID(), UserID: user.GetUserID()})
	}
	if room.IsGroupRoom() {
		room.InformGroupRoom(ctx, user.GetRoomID(), edu_models.OnStudentLeaveGroupRoom, edu_models.NoticeRoom{RoomID: user.GetRoomID(), UserID: user.GetUserID(), UserName: user.GetUserID()})
	}

	service_utils.Push2Client(ctx, param, nil, &edu_models.NoticeRoom{RoomID: p.RoomID, UserID: p.UserID})
	return

}
