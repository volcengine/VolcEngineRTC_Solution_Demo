package edu_service

import (
	"context"
	"encoding/json"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/conn_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/edu_models"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type reconnectParam struct {
	LoginToken string `json:"login_token"`
}

func reconnect(ctx context.Context, param *vc_control.TEventParam) {
	var p reconnectParam
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	conn, err := conn_models.GetConnection(ctx, param.ConnId)
	if err != nil {
		service_utils.Push2ClientWithoutReturn(ctx, param, err)
		return
	}

	err = edu_models.Reconnect(ctx, conn)
	if err != nil {
		logs.Errorf("reconnect failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	user, err := edu_models.GetUser(ctx, conn.GetConnID())
	if err != nil {
		logs.Errorf("get user failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	room, err := edu_models.GetRoom(ctx, user.GetParentRoomID())
	if err != nil {
		logs.Errorf("get room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	//老师直接返回
	if user.GetUserID() == room.GetTeacherUserID() {
		service_utils.Push2Client(ctx, param, nil, edu_models.NoticeRoom{})
		return
	}

	//构造返回结果
	rtcToken, err := genToken(user.GetParentRoomID(), user.GetUserID())
	if err != nil {
		logs.Errorf("gen token failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	interactUsers, err := room.ListInteractStudents(ctx)
	if err != nil {
		logs.Errorf("list interact users failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	res := &joinClassResp{
		RoomInfo:       room.GetRoomInfo(),
		TeacherInfo:    room.GetTeacherInfo(ctx),
		UserInfo:       user.GetUserInfo(),
		Token:          rtcToken,
		CurrentMicUser: interactUsers,
		IsMicOn:        user.IsInteract(),
	}

	if room.IsGroupRoom() {
		groupRoomID := user.GetRoomID()
		idx := room.GetIdxByGroupRoomID(ctx, groupRoomID)
		groupUsers, err := room.ListGroupRoomStudents(ctx, groupRoomID)
		if err != nil {
			logs.Errorf("get group users failed,error:%s", err)
			service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
			return
		}
		//把自己放最上面
		groupUsersRes := make([]*db.EduUserRoomInfo, 0)
		for _, u := range groupUsers {
			if u.UserID == user.GetUserID() {
				groupUsersRes = append([]*db.EduUserRoomInfo{u}, groupUsersRes...)
			} else {
				groupUsersRes = append(groupUsersRes, u)
			}
		}

		GroupRtcToken, err := genToken(user.GetRoomID(), user.GetUserID())
		if err != nil {
			logs.Errorf("gen token failed,error:%s", err)
			service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
			return
		}
		res.GroupRoomID = groupRoomID
		res.RoomIdx = idx
		res.GroupUserList = groupUsersRes
		res.GroupToken = GroupRtcToken
	}

	service_utils.Push2Client(ctx, param, err, res)
}
