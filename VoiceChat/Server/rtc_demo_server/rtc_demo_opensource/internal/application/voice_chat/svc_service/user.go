package svc_service

import (
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_entity"
)

type User struct {
	*svc_entity.SvcUser
	isDirty bool
}

func (u *User) IsDirty() bool {
	return u.isDirty
}

func (u *User) GetRoomID() string {
	return u.RoomID
}

func (u *User) GetUserID() string {
	return u.UserID
}

func (u *User) GetUserName() string {
	return u.UserName
}

func (u *User) GetInteractStatus() int {
	return u.InteractStatus
}

func (u *User) GetSeatID() int {
	return u.SeatID
}

func (u *User) IsEnableInvite() bool {
	return u.NetStatus == svc_db.UserNetStatusOnline && u.InteractStatus == svc_db.UserInteractStatusNormal && u.SeatID == 0
}

func (u *User) IsEnableApply() bool {
	return u.NetStatus == svc_db.UserNetStatusOnline && u.InteractStatus == svc_db.UserInteractStatusNormal && u.SeatID == 0
}

func (u *User) IsEnableInteract() bool {
	return u.NetStatus == svc_db.UserNetStatusOnline && u.InteractStatus == svc_db.UserInteractStatusInviting || u.InteractStatus == svc_db.UserInteractStatusApplying
}

func (u *User) IsInteract() bool {
	return u.InteractStatus == svc_db.UserInteractStatusInteracting
}

func (u *User) IsReconnecting() bool {
	return u.NetStatus == svc_db.UserNetStatusReconnecting
}

func (u *User) IsOnline() bool {
	return u.NetStatus == svc_db.UserNetStatusOnline
}

func (u *User) IsHost() bool {
	return u.UserRole == svc_db.UserRoleHost
}

func (u *User) IsAudience() bool {
	return u.UserRole == svc_db.UserRoleAudience
}

func (u *User) StartLive() {
	u.NetStatus = svc_db.UserNetStatusOnline
	u.InteractStatus = svc_db.UserInteractStatusInteracting
	u.JoinTime = time.Now().UnixNano() / 1e6
	u.isDirty = true
}

func (u *User) JoinRoom(roomID string) {
	u.RoomID = roomID
	u.NetStatus = svc_db.UserNetStatusOnline
	u.JoinTime = time.Now().UnixNano() / 1e6
	u.isDirty = true
}

func (u *User) LeaveRoom() {
	u.NetStatus = svc_db.UserNetStatusOffline
	u.InteractStatus = svc_db.UserInteractStatusNormal
	u.SeatID = 0
	u.LeaveTime = time.Now().UnixNano() / 1e6
	u.isDirty = true
}

func (u *User) SetInteract(interactStatus, seatID int) {
	u.InteractStatus = interactStatus
	u.SeatID = seatID
	u.isDirty = true
}

func (u *User) Disconnect() {
	u.NetStatus = svc_db.UserNetStatusReconnecting
	u.isDirty = true
}

func (u *User) Reconnect(deviceID string) {
	u.NetStatus = svc_db.UserNetStatusOnline
	u.DeviceID = deviceID
	u.isDirty = true
}

func (u *User) SetUpdateTime(time time.Time) {
	u.UpdateTime = time
	u.isDirty = true
}

func (u *User) SetIsDirty(status bool) {
	u.isDirty = status
}

func (u *User) MuteMic() {
	u.Mic = 0
	u.isDirty = true
}

func (u *User) UnmuteMic() {
	u.Mic = 1
	u.isDirty = true
}
