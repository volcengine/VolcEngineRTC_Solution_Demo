package svc_service

import (
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_entity"
)

type Seat struct {
	*svc_entity.SvcSeat
	isDirty bool
}

func (s *Seat) IsDirty() bool {
	return s.isDirty
}

func (s *Seat) GetRoomID() string {
	return s.RoomID
}

func (s *Seat) GetSeatID() int {
	return s.SeatID
}

func (s *Seat) GetOwnerUserID() string {
	return s.OwnerUserID
}

func (s *Seat) IsEnableAlloc() bool {
	return !s.IsLock() && s.OwnerUserID == ""
}

func (s *Seat) IsLock() bool {
	return s.Status == 0
}

func (s *Seat) SetIsDirty(status bool) {
	s.isDirty = status
}

func (s *Seat) Lock() {
	s.Status = 0
	s.isDirty = true
}

func (s *Seat) Unlock() {
	s.Status = 1
	s.isDirty = true
}

func (s *Seat) SetOwnerUserID(userID string) {
	s.OwnerUserID = userID
	s.isDirty = true
}
