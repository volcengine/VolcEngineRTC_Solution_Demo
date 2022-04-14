package cs_entity

import "time"

type CsInteractDetail struct {
	ID           int64     `gorm:"column:id" json:"id"`
	InteractID   string    `gorm:"column:interact_id" json:"interact_id"`
	InteractType int       `gorm:"column:interact_type" json:"interact_type"`
	FromRoomID   string    `gorm:"column:from_room_id" json:"from_room_id"`
	FromUserID   string    `gorm:"column:from_user_id" json:"from_user_id"`
	ToRoomID     string    `gorm:"column:to_room_id" json:"to_room_id"`
	ToUserID     string    `gorm:"column:to_user_id" json:"to_user_id"`
	Status       int       `gorm:"column:status" json:"status"`
	SeatID       int       `gorm:"column:seat_id" json:"seat_id"`
	CreateTime   time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"update_time"`
}

func (id *CsInteractDetail) SetStatus(status int) {
	id.Status = status
}
