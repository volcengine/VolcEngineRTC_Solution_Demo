package svc_service

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_entity"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
)

var seatFactoryClient *SeatFactory

type SeatFactory struct {
	seatRepo SeatRepo
}

func GetSeatFactory() *SeatFactory {
	if seatFactoryClient == nil {
		seatFactoryClient = &SeatFactory{
			seatRepo: GetSeatRepo(),
		}
	}
	return seatFactoryClient
}

func (sf *SeatFactory) NewSeat(ctx context.Context, roomID string, seatID int) *Seat {
	dbSeat := &svc_entity.SvcSeat{
		RoomID: roomID,
		SeatID: seatID,
		Status: 1,
	}
	seat := &Seat{
		SvcSeat: dbSeat,
		isDirty: true,
	}
	return seat
}

func (sf *SeatFactory) Save(ctx context.Context, seat *Seat) error {
	if seat.IsDirty() {
		err := sf.seatRepo.Save(ctx, seat.SvcSeat)
		if err != nil {
			return custom_error.InternalError(err)
		}
		seat.SetIsDirty(false)
	}
	return nil
}

func (sf *SeatFactory) GetSeatByRoomIDSeatID(ctx context.Context, roomID string, seatID int) (*Seat, error) {
	dbSeat, err := sf.seatRepo.GetSeatByRoomIDSeatID(ctx, roomID, seatID)
	if err != nil {
		return nil, custom_error.InternalError(err)
	}
	if dbSeat == nil {
		return nil, nil
	}

	seat := &Seat{
		SvcSeat: dbSeat,
		isDirty: false,
	}
	return seat, nil
}

func (sf *SeatFactory) GetSeatsByRoomID(ctx context.Context, roomID string) ([]*Seat, error) {
	dbSeats, err := sf.seatRepo.GetSeatsByRoomID(ctx, roomID)
	if err != nil {
		return nil, custom_error.InternalError(err)
	}

	seats := make([]*Seat, len(dbSeats))
	for i := 0; i < len(dbSeats); i++ {
		seats[i] = &Seat{
			SvcSeat: dbSeats[i],
			isDirty: false,
		}
	}

	return seats, nil
}
