package public

type EventParam struct {
	AppID     string `json:"app_id"`
	RoomID    string `json:"room_id"`
	UserID    string `json:"user_id"`
	EventName string `json:"event_name"`
	Content   string `json:"content"`
	RequestID string `json:"request_id"`
	DeviceID  string `json:"device_id"`
}
