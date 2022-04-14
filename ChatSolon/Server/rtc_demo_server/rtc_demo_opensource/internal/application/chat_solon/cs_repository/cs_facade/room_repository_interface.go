package cs_facade

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_repository/cs_implement"
)

var roomRepository *cs_implement.RoomRepositoryImpl

func GetRoomRepository() *cs_implement.RoomRepositoryImpl {
	if roomRepository == nil {
		roomRepository = &cs_implement.RoomRepositoryImpl{}
	}
	return roomRepository
}

type RoomRepositoryInterface interface {
	Save(ctx context.Context, room *cs_entity.CsRoom) error

	GetRoomByRoomID(ctx context.Context, roomID string) (*cs_entity.CsRoom, error)
	GetRooms(ctx context.Context) ([]*cs_entity.CsRoom, error)
}
