package cs_entity

import "time"

const (
	RoomStatusPreparing = 0
	RoomStatusStart     = 1
	RoomStatusFinish    = 2
)

type CsRoom struct {
	ID            int64  `gorm:"column:id" json:"id"`
	AppID         string `gorm:"column:app_id" json:"app_id"`
	RoomID        string `gorm:"column:room_id" json:"room_id"`
	RoomName      string `gorm:"column:room_name" json:"room_name"`
	OwnerUserID   string `gorm:"column:owner_user_id" json:"owner_user_id"`
	OwnerUserName string `gorm:"column:owner_user_name" json:"owner_user_name"`
	//Limit string `gorm:"column:limit" json:"limit"`
	Mic        int       `gorm:"column:mic" json:"mic"`
	Camera     int       `gorm:"column:camera" json:"camera"`
	Status     int       `gorm:"column:status" json:"status"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
	UserCount  int       `gorm:"-" json:"user_count"`
	Ext        string    `gorm:"column:ext" json:"ext"`
}

func (r *CsRoom) Prepare() {
	r.Status = RoomStatusPreparing
}

func (r *CsRoom) Start() {
	r.Status = RoomStatusStart
}

func (r *CsRoom) Finish() {
	r.Status = RoomStatusFinish
}
func (r *CsRoom) GetAppID() string {
	return r.AppID
}

func (r *CsRoom) GetRoomID() string {
	return r.RoomID
}
func (r *CsRoom) GetRoomName() string {
	return r.RoomName
}
func (r *CsRoom) GetOwnerUserID() string {
	return r.OwnerUserID
}
func (r *CsRoom) GetOwnerUserName() string {
	return r.OwnerUserName
}

func (r *CsRoom) GetCreateTime() time.Time {
	return r.CreateTime
}

func (r *CsRoom) SetOwnerUserID(userID string) {
	r.OwnerUserID = userID
}
func (r *CsRoom) SetOwnerUserName(userName string) {
	r.OwnerUserName = userName
}
