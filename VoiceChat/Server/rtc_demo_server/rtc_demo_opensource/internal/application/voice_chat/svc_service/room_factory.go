package svc_service

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_db"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_redis"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/public"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"
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

type RoomDetail struct {
	RoomInfo      *Room `json:"room_info"`
	AudienceCount int   `json:"audience_count"`
}

var roomFactoryClient *RoomFactory

type RoomFactory struct {
	roomRepo RoomRepo
}

func GetRoomFactory() *RoomFactory {
	if roomFactoryClient == nil {
		roomFactoryClient = &RoomFactory{
			roomRepo: GetRoomRepo(),
		}
	}
	return roomFactoryClient
}

func (rf *RoomFactory) NewRoom(ctx context.Context, roomName, roomBackgroundImageName, hostUserID, hostUserName string) (*Room, error) {
	roomID, err := ApplyRoomIDWithRetry(ctx)
	if err != nil {
		return nil, err
	}

	dbRoomExt := &svc_entity.SvcRoomExt{
		BackgroundImageName: roomBackgroundImageName,
	}
	ext, _ := json.Marshal(dbRoomExt)
	dbRoom := &svc_entity.SvcRoom{
		AppID:                       config.Configs().SvcAppID,
		RoomID:                      roomID,
		RoomName:                    roomName,
		HostUserID:                  hostUserID,
		HostUserName:                hostUserName,
		Status:                      svc_db.RoomStatusPrepare,
		EnableAudienceInteractApply: 0,
		CreateTime:                  time.Now(),
		UpdateTime:                  time.Now(),
		Ext:                         string(ext),
	}
	room := &Room{
		SvcRoom: dbRoom,
		isDirty: true,
	}
	return room, nil
}

func (rf *RoomFactory) Save(ctx context.Context, room *Room) error {
	if room.IsDirty() {
		room.SetUpdateTime(time.Now())
		err := rf.roomRepo.Save(ctx, room.GetDbRoom())
		if err != nil {
			return custom_error.InternalError(err)
		}
		room.SetIsDirty(false)
	}
	return nil
}

func (rf *RoomFactory) GetRoomByRoomID(ctx context.Context, roomID string) (*Room, error) {
	dbRoom, err := rf.roomRepo.GetRoomByRoomID(ctx, roomID)
	if err != nil {
		return nil, custom_error.InternalError(err)
	}
	if dbRoom == nil {
		return nil, nil
	}

	room := &Room{
		SvcRoom: dbRoom,
		isDirty: true,
	}
	return room, nil
}

func (rf *RoomFactory) GetActiveRoomList(ctx context.Context, needAudienceCount bool) ([]*Room, error) {
	dbRooms, err := rf.roomRepo.GetActiveRooms(ctx)
	if err != nil {
		return nil, custom_error.InternalError(err)
	}
	res := make([]*Room, 0)
	for _, dbRoom := range dbRooms {
		room := &Room{
			SvcRoom: dbRoom,
			isDirty: false,
		}
		res = append(res, room)
	}

	if needAudienceCount {
		roomIDs := make([]string, 0)
		for _, room := range res {
			roomIDs = append(roomIDs, room.GetRoomID())
		}
		roomsAudienceCount, err := rf.GetRoomsAudienceCount(ctx, roomIDs)
		if err != nil {
			logs.CtxError(ctx, "get rooms audience count failed,error:%s", err)
			return nil, err
		}
		for _, room := range res {
			room.AudienceCount = roomsAudienceCount[room.GetRoomID()]
		}
	}

	return res, nil
}

func (rf *RoomFactory) IncrRoomAudienceCount(ctx context.Context, roomID string, count int) error {
	err := svc_redis.IncrRoomAudienceCount(ctx, roomID, int64(count))
	if err != nil {
		return custom_error.InternalError(err)
	}
	return nil
}

func (rf *RoomFactory) GetRoomsAudienceCount(ctx context.Context, roomIDs []string) (map[string]int, error) {
	res := make(map[string]int)
	if len(roomIDs) == 0 {
		return res, nil
	}
	res, err := svc_redis.GetRoomsAudienceCount(ctx, roomIDs)
	if err != nil {
		return nil, custom_error.InternalError(err)
	}
	return res, nil
}

func ApplyRoomIDWithRetry(ctx context.Context) (string, error) {
	roomID, err := generateRoomID(ctx, public.BizIDSvc)
	for i := 0; roomID == 0 && i <= maxRetryNum; i++ {
		backOff := time.Duration(int64(math.Pow(2, float64(i)))) * retryDelay
		if backOff > maxRetryDelay {
			backOff = maxRetryDelay
		}
		time.Sleep(backOff)
		roomID, err = generateRoomID(ctx, public.BizIDSvc)
	}
	if roomID == 0 {
		logs.CtxError(ctx, "failed to generate roomID, err: %s", err)
		return "", custom_error.InternalError(errors.New("make room err"))
	}
	return strconv.FormatInt(roomID, 10), nil
}

func generateRoomID(ctx context.Context, bizID string) (int64, error) {
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

	roomRepo := GetRoomRepo()
	dbRoom, err := roomRepo.GetRoomByRoomID(ctx, strconv.FormatInt(roomID, 10))
	if err == nil && dbRoom == nil {
		return roomID, nil
	}

	if err != nil {
		return 0, custom_error.InternalError(err)
	}

	return 0, custom_error.ErrRoomAlreadyExist
}
