package svc_service

import (
	"context"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_redis"
)

var roomRepoClient RoomRepo
var userRepoClient UserRepo
var seatRepoClient SeatRepo

func GetRoomRepo() RoomRepo {
	if roomRepoClient == nil {
		roomRepoClient = &svc_db.DbRoomRepo{}
	}
	return roomRepoClient
}

func GetUserRepo() UserRepo {
	if userRepoClient == nil {
		userRepoClient = &svc_db.DbUserRepo{}
	}
	return userRepoClient
}

func GetSeatRepo() SeatRepo {
	if seatRepoClient == nil {
		seatRepoClient = &svc_redis.RedisSeatRepo{}
	}
	return seatRepoClient
}

type RoomRepo interface {
	//create or update
	Save(ctx context.Context, room *svc_entity.SvcRoom) error

	//get
	GetRoomByRoomID(ctx context.Context, roomID string) (*svc_entity.SvcRoom, error)
	GetActiveRooms(ctx context.Context) ([]*svc_entity.SvcRoom, error)
}

type UserRepo interface {
	//create or update
	Save(ctx context.Context, user *svc_entity.SvcUser) error

	//update users
	UpdateUsersWithMapByRoomID(ctx context.Context, roomID string, ups map[string]interface{}) error

	//get user
	GetActiveUserByRoomIDUserID(ctx context.Context, roomID, userID string) (*svc_entity.SvcUser, error)
	GetUserByRoomIDUserID(ctx context.Context, roomID, userID string) (*svc_entity.SvcUser, error)
	GetActiveUserByUserID(ctx context.Context, userID string) (*svc_entity.SvcUser, error)

	//get users
	GetActiveUsersByRoomID(ctx context.Context, roomID string) ([]*svc_entity.SvcUser, error)
	GetAudiencesWithoutApplyByRoomID(ctx context.Context, roomID string) ([]*svc_entity.SvcUser, error)
	GetApplyAudiencesByRoomID(ctx context.Context, roomID string) ([]*svc_entity.SvcUser, error)
}

type SeatRepo interface {
	//create or update
	Save(ctx context.Context, seat *svc_entity.SvcSeat) error

	//get
	GetSeatByRoomIDSeatID(ctx context.Context, roomID string, seatID int) (*svc_entity.SvcSeat, error)
	GetSeatsByRoomID(ctx context.Context, roomID string) ([]*svc_entity.SvcSeat, error)
}
