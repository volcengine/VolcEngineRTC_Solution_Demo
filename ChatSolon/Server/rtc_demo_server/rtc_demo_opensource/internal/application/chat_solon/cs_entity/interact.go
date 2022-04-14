package cs_entity

type CsInteract struct {
	ID           int64  `gorm:"column:id" json:"id"`
	InteractID   string `gorm:"interact_id" json:"interact_id"`
	OwnerRoomID  string `gorm:"owner_room_id" json:"owner_room_id"`
	OwnerUserID  string `gorm:"owner_user_id" json:"owner_user_id"`
	RtcAppID     string `gorm:"rtc_app_id" json:"rtc_app_id"`
	RtcRoomID    string `gorm:"rtc_room_id" json:"rtc_room_id"`
	InteractType int    `gorm:"interact_type" json:"interact_type"`
	Status       int    `gorm:"status" json:"status"`
}

func (i *CsInteract) GetInteractID() string {
	return i.InteractID
}
func (i *CsInteract) GetOwnerRoomID() string {
	return i.OwnerRoomID
}
func (i *CsInteract) GetOwnerUserID() string {
	return i.OwnerUserID
}
func (i *CsInteract) GetRtcAppID() string {
	return i.RtcAppID
}
func (i *CsInteract) GetRtcRoomID() string {
	return i.RtcRoomID
}
func (i *CsInteract) GetInteractType() int {
	return i.InteractType
}
func (i *CsInteract) GetStatus() int {
	return i.Status
}

func (i *CsInteract) SetInteractID(interactID string) {
	i.InteractID = interactID
}

func (i *CsInteract) SetOwnerRoomID(ownerRoomID string) {
	i.OwnerRoomID = ownerRoomID
}
func (i *CsInteract) SetOwnerUserID(ownerUserID string) {
	i.OwnerUserID = ownerUserID
}
func (i *CsInteract) SetRtcAppID(rtcAppID string) {
	i.RtcAppID = rtcAppID
}
func (i *CsInteract) SetRtcRoomID(rtcRoomID string) {
	i.RtcRoomID = rtcRoomID
}
func (i *CsInteract) SetType(interactType int) {
	i.InteractType = interactType
}
func (i *CsInteract) SetStatus(status int) {
	i.Status = status
}
