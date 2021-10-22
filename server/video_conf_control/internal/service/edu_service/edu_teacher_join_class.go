package edu_service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/edu_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/conn_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type teacherJoinClassReq struct {
	RoomID     string `json:"room_id"`
	UserName   string `json:"user_name"`
	UserID     string `json:"user_id"`
	LoginToken string `json:"login_token"`
}

func teacherJoinClass(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("teacherJoinClass:%+v", param)
	var p teacherJoinClassReq
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

	conn, err := conn_models.GetConnection(ctx, param.ConnId)
	if err != nil {
		logs.Errorf("failed to get connection, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, err)
		return
	}

	//校验用户是否同时以不同身份在线
	role, err := edu_models.GetUserRoleInOtherRoom(ctx, p.RoomID, p.UserID)
	if err != nil {
		logs.Errorf("get user role in other room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}
	if role != -1 {
		logs.Errorf("teacher in other room,forbid join this class")
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrUserTeacherInOtherClass)
		return
	}

	room, err := edu_models.GetRoom(ctx, p.RoomID)
	if err != nil {
		logs.Errorf("get room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrRoomNotExist)
		return
	}

	if p.UserID != room.GetTeacherUserID() {
		logs.Errorf("user is student of this class,forbid join class")
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrUserRoleNotMatch)
		return
	}

	user, err := edu_models.GetUserByRoomIDUserID(ctx, room.GetRoomID(), p.UserID)
	if err != nil {
		logs.Errorf("get user failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}
	if user == nil {
		user = edu_models.NewUser(&db.EduUserRoomInfo{
			AppID:       config.Config.AppID,
			UserID:      p.UserID,
			UserName:    p.UserName,
			UserRole:    edu_db.UserRoleTeacher,
			CreatedTime: db.EduTime(time.Now()),
			UpdatedTime: db.EduTime(time.Now()),
			IsMicOn:     true,
			IsCameraOn:  true,
			IsHandsUp:   false,
			IsInteract:  false,
			ConnID:      conn.GetConnID(),
			DeviceID:    conn.GetDeviceID(),
		})
	} else {
		//互踢
		user.Inform(ctx, edu_models.OnLogInElsewhere, edu_models.NoticeRoom{})
		user.SetName(p.UserName)
		user.SetConnID(conn.GetConnID())
		user.SetDeviceID(conn.GetDeviceID())
	}

	user.TeacherJoinRoom(ctx, room)
	token, err := genToken(p.RoomID, p.UserID)
	if err != nil {
		logs.Errorf("gen token failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}
	user.SetToken(token)

	err = user.Save(ctx)
	if err != nil {
		logs.Errorf("teacher join room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	room.InformRoom(ctx, edu_models.OnTeacherJoinClass, &edu_models.NoticeRoom{RoomID: room.GetRoomID()})

	service_utils.Push2Client(ctx, param, err, user.GetUserInfo())

}
