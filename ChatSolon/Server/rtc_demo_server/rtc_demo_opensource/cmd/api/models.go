package api

const (
	EventTypeUserLeaveRoom = "UserLeaveRoom"
)

const (
	LeaveRoomReasonUserLeave      = "userLeave"
	LeaveRoomReasonConnectionLost = "connectionLost"
)

type RtmParam struct {
	Message string `json:"message"`
}

type RtmCallbackParam struct {
	EventType string `json:"EventType"`
	EventData string `json:"EventData"`
	EventTime string `json:"EventTime"`
	EventId   string `json:"EventId"`
	AppId     string `json:"AppId"`
	Version   string `json:"Version"`
	Signature string `json:"Signature"`
	Nonce     string `json:"Nonce"`
}

type EventDataLeaveRoom struct {
	RoomId     string `json:"RoomId"`
	UserId     string `json:"UserId"`
	DeviceType string `json:"DeviceType"`
	Reason     string `json:"Reason"`
	Timestamp  int64  `json:"Timestamp"`
}
