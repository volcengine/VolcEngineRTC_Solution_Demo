package svc_redis

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/voice_chat/svc_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/redis_cli"
)

const (
	keyRoomSeatPrefix = "svc:vccontrol:room_seat:"
)

type RedisSeatRepo struct{}

func (repo *RedisSeatRepo) Save(ctx context.Context, seat *svc_entity.SvcSeat) error {
	data, _ := json.Marshal(seat)
	err := redis_cli.Client.HSet(ctx, roomSeatKey(seat.RoomID), strconv.Itoa(seat.SeatID), data).Err()
	if err != nil {
		return err
	}
	redis_cli.Client.Expire(ctx, roomSeatKey(seat.RoomID), expireTime)
	return nil
}

func (repo *RedisSeatRepo) GetSeatByRoomIDSeatID(ctx context.Context, roomID string, seatID int) (*svc_entity.SvcSeat, error) {
	data, err := redis_cli.Client.HGet(ctx, roomSeatKey(roomID), strconv.Itoa(seatID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	rs := &svc_entity.SvcSeat{}
	err = json.Unmarshal([]byte(data), rs)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (repo *RedisSeatRepo) GetSeatsByRoomID(ctx context.Context, roomID string) ([]*svc_entity.SvcSeat, error) {
	dataMap, err := redis_cli.Client.HGetAll(ctx, roomSeatKey(roomID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	res := make([]*svc_entity.SvcSeat, 0)
	for _, data := range dataMap {
		rs := &svc_entity.SvcSeat{}
		err = json.Unmarshal([]byte(data), rs)
		if err != nil {
			return nil, err
		}
		res = append(res, rs)
	}
	return res, nil
}

func roomSeatKey(roomID string) string {
	return keyRoomSeatPrefix + roomID
}
