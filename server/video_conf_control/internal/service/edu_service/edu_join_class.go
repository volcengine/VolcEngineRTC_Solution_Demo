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
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/pkg/token"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/service/service_utils"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/kitex_gen/vc_control"
)

type joinClassReq struct {
	UserID      string `json:"user_id"`
	RoomID      string `json:"room_id"`
	LoginToken  string `json:"login_token"`
	UserName    string `json:"user_name"`
	IsReconnect bool   `json:"is_reconnect"`
}

type joinClassResp struct {
	Token          string                `json:"token"`
	RoomInfo       *db.EduRoomInfo       `json:"room_info"`
	TeacherInfo    *db.EduUserRoomInfo   `json:"teacher_info"`
	UserInfo       *db.EduUserRoomInfo   `json:"user_info"`
	CurrentMicUser []*db.EduUserRoomInfo `json:"current_mic_user"`
	GroupToken     string                `json:"group_token"`
	GroupUserList  []*db.EduUserRoomInfo `json:"group_user_list"`
	RoomIdx        int                   `json:"room_idx"`
	GroupRoomID    string                `json:"group_room_id"`
	IsMicOn        bool                  `json:"is_mic_on"`
}

func joinClass(ctx context.Context, param *vc_control.TEventParam) {
	logs.Infof("joinClass:%+v", param)
	var p joinClassReq
	if err := json.Unmarshal([]byte(param.Content), &p); err != nil {
		logs.Warnf("input format error, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//校验参数
	if p.UserID == "" || p.RoomID == "" || p.UserName == "" {
		logs.Warnf("input user_id 、user_name or room_id error, params: %v", p)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrInput)
		return
	}

	//校验用户是否同时以不同身份在其它房间
	role, err := edu_models.GetUserRoleInOtherRoom(ctx, p.RoomID, p.UserID)
	if err != nil {
		logs.Errorf("get user role in other room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}
	if role == edu_db.UserRoleTeacher {
		logs.Errorf("user is teacher role in other room,forbid join this class")
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrUserTeacherInOtherClass)
		return
	}

	room, err := edu_models.GetRoom(ctx, p.RoomID)
	if err != nil {
		logs.Errorf("get room failed,roomid:%s,error:%s", p.RoomID, err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrRoomNotExist)
		return
	}

	if p.UserID == room.GetTeacherUserID() {
		logs.Errorf("user is teacher of this class,forbid join class")
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.ErrUserRoleNotMatch)
		return
	}

	conn, err := conn_models.GetConnection(ctx, param.ConnId)
	if err != nil {
		logs.Errorf("failed to get connection, err: %v", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, err)
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
			UserRole:    edu_db.UserRoleStudent,
			CreatedTime: db.EduTime(time.Now()),
			UpdatedTime: db.EduTime(time.Now()),
			IsMicOn:     false,
			IsCameraOn:  false,
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

	err = user.JoinRoom(ctx, room)
	if err != nil {
		logs.Errorf("join room failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}
	logs.Infof("join room success,user:%v,room:%v", user, room)

	//持久化
	err = user.Save(ctx)
	if err != nil {
		logs.Errorf("save user failed,error:%s", err)
		service_utils.Push2ClientWithoutReturn(ctx, param, custom_error.InternalError(err))
		return
	}

	//通知小组成员
	if room.IsGroupRoom() {
		room.InformGroupRoom(ctx, user.GetRoomID(), edu_models.OnStudentJoinGroupRoom, edu_models.NoticeRoom{RoomID: user.GetRoomID(), UserID: user.GetUserID(), UserName: user.GetUserName()})
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

func genToken(roomID, userID string) (string, error) {
	//生成token
	tokenParam := &token.GenerateParam{
		AppID:        config.Config.AppID,
		AppKey:       config.Config.AppKey,
		RoomID:       roomID,
		UserID:       userID,
		ExpireAt:     7 * 24 * 3600,
		CanPublish:   true,
		CanSubscribe: true,
		CanSwitch:    true,
	}
	return token.GenerateToken(tokenParam)
}

//todo 自定义内部错误标识
