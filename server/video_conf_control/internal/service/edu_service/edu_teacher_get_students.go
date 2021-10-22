package edu_service

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/edu_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type teacherGetStudentsInfoReq struct {
	RoomID     string `json:"room_id"`
	PageNumber int    `json:"page_number"`
	PageSize   int    `json:"page_size"`
	LoginToken string `json:"login_token"`
}

type teacherGetStudentsInfoResp struct {
	RoomID        string                        `json:"room_id"`
	RoomType      int                           `json:"room_type"`
	UserList      []*db.EduUserRoomInfo         `json:"user_list"`
	GroupUserList map[int][]*db.EduUserRoomInfo `json:"group_user_list"`
	StudentCount  int                           `json:"student_count"`
}

func teacherGetStudentsInfo(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("teacherJoinClass:%+v", param)
	var p teacherGetStudentsInfoReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//校验参数
	if p.RoomID == "" {
		logs.Warnf("input group_room_id error, params: %v", p)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	if p.PageNumber <= 0 {
		p.PageNumber = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}

	room, err := edu_models.GetRoom(ctx, p.RoomID)
	if err != nil {
		logs.Errorf("get room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrRoomNotExist)
		return
	}

	studentCount, err := room.GetRoomStudentCount(ctx)
	if err != nil {
		logs.Errorf("get student count failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}
	res := &teacherGetStudentsInfoResp{
		RoomID:       p.RoomID,
		RoomType:     room.GetRoomType(),
		StudentCount: studentCount,
	}

	if p.PageNumber <= 0 {
		p.PageNumber = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if room.GetRoomType() == edu_db.RoomTypeSingleRoomClass {
		users, err := room.ListStudents(ctx, p.PageNumber, p.PageSize)
		if err != nil {
			logs.Errorf("get user list failed,error:%s", err)
			service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
			return
		}
		res.UserList = users
	} else {
		groupUsers, err := room.ListStudentsByGroup(ctx, 1, 35)
		if err != nil {
			logs.Errorf("get group user list failed,error:%s", err)
			service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
			return
		}
		res.GroupUserList = groupUsers
	}
	service_utils.Push2Client(ctx, param, err, res)
}
