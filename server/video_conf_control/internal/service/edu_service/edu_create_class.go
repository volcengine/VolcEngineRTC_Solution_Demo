package edu_service

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/edu_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

const (
	maxUserNameLen = 100
	maxRoomNameLen = 200
)

type eduCreateClassReq struct {
	AppID          string `json:"app_id"`
	RoomID         string `json:"room_id"`
	RoomName       string `json:"room_name"`
	TeacherName    string `json:"teacher_name"`
	RoomType       int    `json:"room_type"`
	CreateUserID   string `json:"create_user_id"`
	BeginClassTime int64  `json:"begin_class_time"`
	EndClassTime   int64  `json:"end_class_time"`
	GroupNum       int    `json:"group_num"`
	GroupLimit     int    `json:"group_limit"`
	LoginToken     string `json:"login_token"`
	Token          string `json:"token"`
}

func eduCreateClass(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("eduCreateClass:%+v", param)
	var p eduCreateClassReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//校验参数
	if p.CreateUserID == "" || p.RoomName == "" {
		logs.Warnf("input format error, params: %v", p)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}
	//url解码
	roomName, err := url.QueryUnescape(p.RoomName)
	if err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	if len(p.CreateUserID) > maxUserNameLen || len(roomName) > maxRoomNameLen {
		logs.Warnf("input format error, params: %v", p)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	p.AppID = config.Config.AppID

	//校验用户是否同时以不同身份在线
	role, err := edu_models.GetUserRoleInOtherRoom(ctx, p.RoomID, p.CreateUserID)
	if err != nil {
		logs.Errorf("get user role in other room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}
	if role != -1 {
		logs.Errorf("teacher in other room,forbid join this class")
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrUserStudentInOtherClass)
		return
	}

	//创建房间
	roomID, err := edu_models.ApplyRoomIDWithRetry(ctx)
	if err != nil {
		logs.Errorf("apply room id failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}
	room := edu_models.NewRoom(&db.EduRoomInfo{
		AppID:             config.Config.AppID,
		RoomID:            roomID,
		RoomName:          roomName,
		RoomType:          p.RoomType,
		CreateUserID:      p.CreateUserID,
		Status:            edu_db.ClassPending,
		BeginClassTime:    p.BeginClassTime,
		EndClassTime:      p.EndClassTime,
		AudioMuteAll:      true,
		VideoMuteAll:      true,
		EnableGroupSpeech: false,
		EnableInteractive: false,
		IsRecording:       false,
		CreatedTime:       db.EduTime(time.Now()),
		UpdatedTime:       db.EduTime(time.Now()),
		TeacherName:       p.TeacherName,
		GroupNum:          config.Config.EduGroupNum,
		GroupLimit:        config.Config.EduGroupLimit,
	})

	err = room.Save(ctx)
	if err != nil {
		logs.Errorf("save room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	res := room.GetRoomInfo()
	service_utils.Push2Client(ctx, param, err, res)
}
