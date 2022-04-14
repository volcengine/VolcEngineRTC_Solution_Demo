package login_entity

import "time"

type UserProfile struct {
	ID        int64     `gorm:"column:id" json:"id"`
	UserID    string    `gorm:"column:user_id" json:"user_id"`
	UserName  string    `gorm:"column:user_name" json:"user_name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
