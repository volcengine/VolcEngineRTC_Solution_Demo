package rtc_openapi

type sendUnicastParam struct {
	AppID   string `json:"AppId"`
	From    string `json:"From"`
	To      string `json:"To"`
	Binary  bool   `json:"Binary"`
	Message string `json:"Message"`
}

type sendRoomUnicastParam struct {
	AppID   string `json:"AppId"`
	RoomID  string `json:"RoomId"`
	From    string `json:"From"`
	To      string `json:"To"`
	Binary  bool   `json:"Binary"`
	Message string `json:"Message"`
}

type sendBroadcastParam struct {
	AppID   string `json:"AppId"`
	RoomID  string `json:"RoomId"`
	From    string `json:"From"`
	Binary  bool   `json:"Binary"`
	Message string `json:"Message"`
}
