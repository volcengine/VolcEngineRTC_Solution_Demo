package login_entity

import "time"

type LoginToken struct {
	Token      string    `gorm:"column:token" json:"token"`
	UserID     string    `gorm:"column:user_id" json:"user_id"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}
