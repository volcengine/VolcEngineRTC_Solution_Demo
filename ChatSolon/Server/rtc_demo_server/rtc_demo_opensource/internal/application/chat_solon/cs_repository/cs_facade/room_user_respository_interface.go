package cs_facade

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_repository/cs_implement"
)

var roomUserRepository *cs_implement.RoomUserRepositoryImpl

func GetRoomUserRepository() *cs_implement.RoomUserRepositoryImpl {
	if roomUserRepository == nil {
		roomUserRepository = &cs_implement.RoomUserRepositoryImpl{}
	}
	return roomUserRepository
}

type RoomUserRepositoryInterface interface {
	Save(ctx context.Context, user *cs_entity.CsRoomUser) error

	GetActiveUserByRoomIDUserID(ctx context.Context, roomID, userID string) (*cs_entity.CsRoomUser, error)
	GetReconnectUserByRoomIDUserID(ctx context.Context, roomID, userID string) (*cs_entity.CsRoomUser, error)
	GetActiveUsersByRoomID(ctx context.Context, roomID string) ([]*cs_entity.CsRoomUser, error)
	GetActiveUsersByRoomIDStatus(ctx context.Context, roomID string, status []int) ([]*cs_entity.CsRoomUser, error)
	GetActiveAudiencesByRoomIDStatus(ctx context.Context, roomID string, status []int) ([]*cs_entity.CsRoomUser, error)
	GetActiveUsersByRoomIDRole(ctx context.Context, roomID string, role int) ([]*cs_entity.CsRoomUser, error)
	GetUserByStatusOrderUtime(ctx context.Context, roomID string, status int) (*cs_entity.CsRoomUser, error)

	UpdateUsersByUserID(ctx context.Context, userID string, ups map[string]interface{}) error
	UpdateUsersByRoomID(ctx context.Context, roomID string, ups map[string]interface{}) error

	GetUserCountByRoomID(ctx context.Context, roomID string) (int, error)
}
