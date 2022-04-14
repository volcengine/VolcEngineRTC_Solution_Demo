package cs_service

import (
	"context"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/chat_solon/cs_repository/cs_facade"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/redis_cli/general"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/redis_cli/lock"
)

const (
	retryDelay    = 8 * time.Millisecond
	maxRetryDelay = 128 * time.Millisecond
	maxRetryNum   = 10
)

const (
	roomLength = 8
)

type Hall struct {
	roomRepository     cs_facade.RoomRepositoryInterface
	roomUserRepository cs_facade.RoomUserRepositoryInterface
}

var hall *Hall

func GetHall() *Hall {
	if hall == nil {
		hall = &Hall{
			roomRepository:     cs_facade.GetRoomRepository(),
			roomUserRepository: cs_facade.GetRoomUserRepository(),
		}
	}
	return hall
}

type CreateRoomParam struct {
	AppID        string
	RoomID       string
	RoomName     string
	HostUserID   string
	HostUserName string
}

func (h *Hall) CreateRoom(ctx context.Context, param *CreateRoomParam) (*cs_entity.CsRoom, *cs_entity.CsRoomUser, error) {
	room := &cs_entity.CsRoom{
		AppID:         param.AppID,
		RoomID:        param.RoomID,
		RoomName:      param.RoomName,
		OwnerUserID:   param.HostUserID,
		OwnerUserName: param.HostUserName,
		Status:        cs_entity.RoomStatusPreparing,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		UserCount:     1,
	}
	err := h.roomRepository.Save(ctx, room)
	if err != nil {
		logs.CtxError(ctx, "save room failed,error:%s", err)
		return nil, nil, custom_error.InternalError(err)
	}
	host := &cs_entity.CsRoomUser{
		AppID:      param.AppID,
		RoomID:     param.RoomID,
		UserID:     param.HostUserID,
		UserName:   param.HostUserName,
		UserRole:   cs_entity.UserRoleHost,
		Mic:        1,
		Camera:     1,
		NetStatus:  cs_entity.UserNetStatusOnline,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	err = h.roomUserRepository.Save(ctx, host)
	if err != nil {
		logs.CtxError(ctx, "save user failed,error:%s", err)
		return nil, nil, custom_error.InternalError(err)
	}
	return room, host, nil
}

func (h *Hall) CleanUserResidual(ctx context.Context, userID string) error {
	err := h.roomUserRepository.UpdateUsersByUserID(ctx, userID, map[string]interface{}{
		"net_status": cs_entity.UserNetStatusOffline,
	})
	if err != nil {
		return custom_error.InternalError(err)
	}
	return nil
}

func (h *Hall) ListRooms(ctx context.Context) ([]*cs_entity.CsRoom, error) {
	rooms, err := h.roomRepository.GetRooms(ctx)
	if err != nil {
		logs.CtxError(ctx, "get rooms failed,error:%s", err)
		return nil, custom_error.InternalError(err)
	}
	return rooms, nil
}

func (h *Hall) GenerateRoomIDWithRetry(ctx context.Context) (string, error) {
	roomID, err := h.generateRoomID(ctx, public.BizIDCs)
	for i := 0; roomID == 0 && i <= maxRetryNum; i++ {
		backOff := time.Duration(int64(math.Pow(2, float64(i)))) * retryDelay
		if backOff > maxRetryDelay {
			backOff = maxRetryDelay
		}
		time.Sleep(backOff)
		roomID, err = h.generateRoomID(ctx, public.BizIDCs)
	}
	if roomID == 0 {
		logs.CtxError(ctx, "failed to generate roomID, err: %s", err)
		return "", custom_error.InternalError(errors.New("make room err"))
	}
	return strconv.FormatInt(roomID, 10), nil
}

func (h *Hall) generateRoomID(ctx context.Context, bizID string) (int64, error) {
	ok, lt := lock.LockGenerateRoomID(ctx)
	if !ok {
		return 0, custom_error.ErrLockRedis
	}
	defer lock.UnlockGenerateRoomID(ctx, lt)

	roomID, err := general.GetGeneratedRoomID(ctx, bizID)
	if err != nil {
		return 0, custom_error.InternalError(err)
	}

	baseline := int64(math.Pow10(roomLength))
	minRoomID := int64(math.Pow10(roomLength - 1))

	if roomID == 0 {
		roomID = time.Now().Unix() % baseline
	} else {
		roomID = (roomID + 1) % baseline
	}

	if roomID < minRoomID {
		roomID += minRoomID
	}

	general.SetGeneratedRoomID(ctx, bizID, roomID)

	room, err := h.roomRepository.GetRoomByRoomID(ctx, strconv.FormatInt(roomID, 10))
	if err == nil && room == nil {
		return roomID, nil
	}

	if err != nil {
		return 0, custom_error.InternalError(err)
	}

	return 0, custom_error.ErrRoomAlreadyExist
}
