package svc_entity

type SvcSeat struct {
	RoomID      string `gorm:"column:room_id" json:"room_id"`
	SeatID      int    `gorm:"column:seat_id" json:"seat_id"` //1-8
	OwnerUserID string `gorm:"column:user_id" json:"user_id"` //申请中的不算，真正上麦的userID
	Status      int    `gorm:"column:status" json:"status"`   //0:lock  1:unlock
}
