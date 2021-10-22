package edu_service

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type teacherGetStuMicOnListReq struct {
	LoginToken string `json:"login_token"`
	RoomID     string `json:"room_id"`
}

type teacherGetStuMicOnListResp struct {
	RoomID   string                `json:"room_id"`
	UserList []*db.EduUserRoomInfo `json:"user_list"`
}

func teacherGetStuMicOnList(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("teacherGetStudentOnMicList:%+v", param)
	var p teacherGetStuMicOnListReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//校验参数
	if p.RoomID == "" {
		logs.Warnf("input room_id error, params: %v", p)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	room, err := edu_models.GetRoom(ctx, p.RoomID)
	if err != nil {
		logs.Errorf("get room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrRoomNotExist)
		return
	}

	users, err := room.ListInteractStudents(ctx)
	if err != nil {
		logs.Errorf("list hands up users failed,error: %s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	res := &teacherGetStuMicOnListResp{
		RoomID:   p.RoomID,
		UserList: users,
	}
	service_utils.Push2Client(ctx, param, err, res)
}
