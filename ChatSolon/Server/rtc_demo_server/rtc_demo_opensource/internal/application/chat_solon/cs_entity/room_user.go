package cs_entity

import "time"

const (
	UserRoleHost     = 0
	UserRoleAudience = 1
)

const (
	UserNetStatusOffline      = 0
	UserNetStatusOnline       = 1
	UserNetStatusReconnecting = 2
)

type CsRoomUser struct {
	ID             int64     `gorm:"column:id" json:"id"`
	AppID          string    `gorm:"column:app_id" json:"app_id"`
	RoomID         string    `gorm:"column:room_id" json:"room_id"`
	UserID         string    `gorm:"column:user_id" json:"user_id"`
	UserName       string    `gorm:"column:user_name" json:"user_name"`
	UserRole       int       `gorm:"column:user_role" json:"user_role"`
	Mic            int       `gorm:"column:mic" json:"mic"`
	Camera         int       `gorm:"column:camera" json:"camera"`
	NetStatus      int       `gorm:"column:net_status" json:"net_status"`
	InteractStatus int       `gorm:"column:interact_status" json:"interact_status"`
	CreateTime     time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime     time.Time `gorm:"column:update_time" json:"update_time"`
	DeviceID       string    `gorm:"column:device_id" json:"device_id"`
}

func (u *CsRoomUser) IsHost() bool {
	return u.UserRole == UserRoleHost
}

func (u *CsRoomUser) IsAudience() bool {
	return u.UserRole == UserRoleAudience
}

func (u *CsRoomUser) GetNetStatus() int {
	return u.NetStatus
}

func (u *CsRoomUser) GetInteractStatus() int {
	return u.InteractStatus
}

func (u *CsRoomUser) GetUserID() string {
	return u.UserID
}

func (u *CsRoomUser) GetUserName() string {
	return u.UserName
}

func (u *CsRoomUser) SetInteractStatus(status int) {
	u.InteractStatus = status
}

func (u *CsRoomUser) Mute() {
	u.Mic = 0
}
func (u *CsRoomUser) Unmute() {
	u.Mic = 1
}

func (u *CsRoomUser) Leave() {
	u.NetStatus = UserNetStatusOffline
}
