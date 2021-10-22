package edu_models

import (
	"context"
	"errors"
	"time"

	logs "github.com/sirupsen/logrus"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/db/edu_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/models/conn_models"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/pkg/record"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/pkg/response"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/rpc/frontier"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/pkg/public"
)

const (
	MaxReconnectRetry = 3
)

type User struct {
	u       *db.EduUserRoomInfo
	isDirty bool
	//use struct for extern attribution
}

//-1表示不存在
func GetUserRoleInOtherRoom(ctx context.Context, roomID, userID string) (int, error) {
	notExist := -1
	users, err := UserDsClient.GetActiveUsersByUserID(ctx, userID)
	if err != nil {
		return notExist, err
	}

	for _, user := range users {
		if user.ParentRoomID != roomID {
			return user.UserRole, nil
		}
	}
	return notExist, nil
}

func NewUser(dbUser *db.EduUserRoomInfo) *User {
	user := &User{
		u:       dbUser,
		isDirty: true,
	}
	return user

}

func GetUser(ctx context.Context, connID string) (*User, error) {
	dbUser, err := UserDsClient.GetActiveUserByConnID(ctx, connID)
	if err != nil {
		return nil, err
	}
	if dbUser == nil {
		return nil, nil
	}

	user := &User{
		u:       dbUser,
		isDirty: false,
	}
	return user, nil
}

func GetUserByRoomIDUserID(ctx context.Context, roomID, userID string) (*User, error) {
	dbUser, err := UserDsClient.GetUserByParentRoomIDUserID(ctx, roomID, userID)
	if err != nil {
		return nil, err
	}
	if dbUser == nil || dbUser.UserStatus != edu_db.UserStatusOnline {
		return nil, nil
	}

	user := &User{
		u:       dbUser,
		isDirty: false,
	}
	return user, nil

}

func (u *User) Save(ctx context.Context) error {
	if u.isDirty {
		u.u.UpdatedTime = db.EduTime(time.Now())
		err := UserDsClient.Save(ctx, u.u)
		if err != nil {
			return err
		}
	}
	u.isDirty = false
	return nil
}

//teacher
func (u *User) TeacherJoinRoom(ctx context.Context, room *Room) {
	u.u.UserStatus = edu_db.UserStatusOnline
	u.u.JoinTime = time.Now().UnixNano()

	u.u.ParentRoomID = room.GetRoomID()
	u.u.RoomID = room.GetRoomID()

	u.isDirty = true

}

func (u *User) TeacherLeaveRoom(ctx context.Context) {
	u.u.UserStatus = edu_db.UserStatusOffline
	u.u.LeaveTime = time.Now().UnixNano()

	u.u.IsHandsUp = false
	u.u.IsInteract = false
	u.u.IsMicOn = true
	u.u.IsCameraOn = true

	u.isDirty = true
}

//student
func (u *User) JoinRoom(ctx context.Context, room *Room) error {
	//更新字段
	u.u.UserStatus = edu_db.UserStatusOnline
	u.u.JoinTime = time.Now().UnixNano()

	//没有roomID
	if u.u.ParentRoomID == "" && u.u.RoomID == "" {
		u.u.ParentRoomID = room.GetRoomID()
		if room.IsGroupRoom() {
			//已退房用户保留group room id
			dbUser, err := UserDsClient.GetUserByParentRoomIDUserID(ctx, room.GetRoomID(), u.u.UserID)
			if err != nil {
				logs.Errorf("get user failed,error:%s", err)
				return err
			}
			if dbUser != nil {
				u.u.RoomID = dbUser.RoomID
			} else {
				//第一次进房
				groupRoomID, err := room.GetIdleGroupRoomID(ctx)
				if err != nil {
					return err
				}
				u.u.RoomID = groupRoomID
			}

		} else {
			u.u.RoomID = room.GetRoomID()
		}
	}

	u.isDirty = true
	return nil
}

func (u *User) LeaveRoom(ctx context.Context) {
	u.u.UserStatus = edu_db.UserStatusOffline
	u.u.LeaveTime = time.Now().UnixNano()

	u.u.IsHandsUp = false
	u.u.IsInteract = false
	u.u.IsMicOn = true
	u.u.IsCameraOn = true

	u.isDirty = true
}

func (u *User) Inform(ctx context.Context, event InformEvent, data interface{}) {
	defer public.CheckPanic()

	c, err := conn_models.GetConnection(ctx, u.u.ConnID)
	if err != nil {
		logs.Errorf("failed to get connection, err: %v", err)
		return
	}
	frontier.PushToClient(ctx, c.GetConnID(), c.GetAddr(), c.GetAddr6(), string(event), response.NewInformToClient(string(event), data), -1)
}

func (u *User) Disconnect(ctx context.Context) error {
	u.u.UserStatus = edu_db.UserStatusReconnecting
	u.isDirty = true
	err := u.Save(ctx)
	if err != nil {
		logs.Errorf("save user failed,user:%#v,error:%s", u.u, err)
	}

	go func() {
		time.Sleep(time.Duration(config.Config.ReconnectTimeout) * time.Second)

		dbUser, err := UserDsClient.GetUserByParentRoomIDUserID(ctx, u.u.ParentRoomID, u.u.UserID)
		if err != nil {
			logs.Errorf("get user failed,error:%s", err)
			return
		}
		logs.Infof("user reconnecting : %#v", dbUser)
		if dbUser != nil && dbUser.UserStatus == edu_db.UserStatusReconnecting {
			user := NewUser(dbUser)
			room, err := GetRoom(ctx, user.GetParentRoomID())
			if err != nil {
				logs.Errorf("get room failed,error:%s", err)
				return
			}

			if user.IsTeacher() {
				user.TeacherLeaveRoom(ctx)
				err = user.Save(ctx)
				if err != nil {
					logs.Errorf("save user failed,error:%s", err)
					return
				}
				room.InformRoom(ctx, OnTeacherLeaveClass, NoticeRoom{RoomID: user.GetParentRoomID(), UserID: user.GetUserID(), UserName: user.GetUserName()})
			} else {
				user.LeaveRoom(ctx)
				err = user.Save(ctx)
				if err != nil {
					logs.Errorf("save user failed,error:%s", err)
					return
				}
				if room.IsGroupRoom() {
					room.InformGroupRoom(ctx, user.GetRoomID(), OnStudentLeaveGroupRoom, NoticeRoom{RoomID: user.GetRoomID(), UserID: user.GetUserID(), UserName: user.GetUserID()})
				}
			}
		}
	}()

	return nil
}

func Reconnect(ctx context.Context, c *conn_models.Connection) error {
	dbUser, err := getReconnectingUser(ctx, c)
	if err != nil {
		logs.Errorf("get reconnect user failed,error:%s", err)
		return err
	}
	if dbUser == nil {
		logs.Warnf("get reconnect user failed,user not exist")
		return errors.New("reconnect user not exist")
	}

	dbUser.UserStatus = edu_db.UserStatusOnline
	dbUser.ConnID = c.GetConnID()
	err = UserDsClient.Save(ctx, dbUser)
	if err != nil {
		logs.Errorf("reconnect failed,user:%v,error:%s", dbUser, err)
		return err
	}
	return nil
}

func getReconnectingUser(ctx context.Context, c *conn_models.Connection) (user *db.EduUserRoomInfo, err error) {
	deviceID := c.GetDeviceID()
	if deviceID == "" {
		return nil, errors.New("device id is empty")
	}
	for i := 0; i < MaxReconnectRetry; i++ {
		user, err = UserDsClient.GetReconnectUserByDeviceID(ctx, c.GetDeviceID())
		if err == nil && user != nil {
			return user, nil
		}
		logs.Warnf("cannot find reconnecting user, wait 10ms")
		time.Sleep(10 * time.Millisecond)
	}

	return nil, err
}

func (u *User) GetParentRoomID() string {
	return u.u.ParentRoomID
}

func (u *User) GetRoomID() string {
	return u.u.RoomID
}

func (u *User) GetUserID() string {
	return u.u.UserID
}
func (u *User) GetUserName() string {
	return u.u.UserName
}

func (u *User) SetConnID(connID string) {
	u.u.ConnID = connID
	u.isDirty = true
}

func (u *User) GetConnID() string {
	return u.u.ConnID
}

func (u *User) GetUserInfo() *db.EduUserRoomInfo {
	return u.u
}

func (u *User) IsInteract() bool {
	return u.u.IsInteract
}

func (u *User) IsHandsUp() bool {
	return u.u.IsHandsUp
}

func (u *User) StartInteract() {
	u.u.IsInteract = true
	u.u.IsHandsUp = false
	u.isDirty = true

}

func (u *User) FinishInteract() {
	u.u.IsInteract = false
	u.isDirty = true
}

func (u *User) HandsUp() {
	u.u.IsHandsUp = true
	u.isDirty = true
}

func (u *User) CancelHandsUp() {
	u.u.IsHandsUp = false
	u.isDirty = true
}

func (u *User) TurnOnCamera() {
	u.u.IsCameraOn = true
	u.isDirty = true
}

func (u *User) TurnOffCamera() {
	u.u.IsCameraOn = false
	u.isDirty = true

}

func (u *User) TurnOnMic() {
	u.u.IsMicOn = true
	u.isDirty = true
}
func (u *User) TurnOffMic() {
	u.u.IsMicOn = false
	u.isDirty = true
}

func (u *User) SetName(userName string) {
	u.u.UserName = userName
	u.isDirty = true
}
func (u *User) SetDeviceID(deviceID string) {
	u.u.DeviceID = deviceID
	u.isDirty = true

}

func (u *User) SetToken(token string) {
	u.u.RtcToken = token
	u.isDirty = true
}

func (u *User) GetRole() int {
	return u.u.UserRole
}

func (u *User) IsTeacher() bool {
	return u.u.UserRole == edu_db.UserRoleTeacher
}

func (u *User) IsStudent() bool {
	return u.u.UserRole == edu_db.UserRoleStudent
}

func StartRecord(ctx context.Context, appID, roomID, userID, roomName string) {
	if taskID, err := record.StartRecord(ctx, []string{userID}, "", appID, roomID, record.SingleStreamMode); err != nil {
		logs.Errorf("start record failed,roomID:%s,userID:%s,error:%s", roomID, userID, err)
	} else {
		if err := edu_db.StartRecord(ctx, appID, roomID, userID, taskID, roomName); err != nil {
			logs.Errorf("db start record failed,taskID:%s,error:%s", taskID, err)
		}
	}
}

func StopRecord(ctx context.Context, appID, roomID, userID string) {
	if taskID, err := edu_db.QueryTaskId(ctx, appID, roomID, userID); err != nil {
		logs.Errorf("get taskID failed,err:%s", err)
	} else {
		if err := record.StopRecord(ctx, appID, roomID, taskID); err != nil {
			logs.Errorf("stop record failed,err:%s", err)
		}
	}
}
