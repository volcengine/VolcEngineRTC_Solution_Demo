package svc_redis

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/redis_cli"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
)

const expireTime = 12 * time.Hour

const (
	keyRoomAudienceCountPrefix = "svc:vccontrol:room_audience_count:"
)

func GetRoomsAudienceCount(ctx context.Context, roomIDs []string) (map[string]int, error) {
	pl := redis_cli.Client.Pipeline()
	defer func() {
		_ = pl.Close()
	}()

	cmdMap := make(map[string]*redis.StringCmd)
	for _, roomID := range roomIDs {
		cmdMap[roomID] = pl.Get(ctx, roomAudienceCountKey(roomID))
	}

	_, err := pl.Exec(ctx)
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}

	}

	res := make(map[string]int)
	for roomID, cmd := range cmdMap {
		countString, err := cmd.Result()
		logs.CtxInfo(ctx, "redis cmd result:%s,err:%s", countString, err)
		if err == redis.Nil {
			res[roomID] = 0
		} else {
			count, _ := strconv.Atoi(countString)
			res[roomID] = count
		}
	}

	return res, nil
}

func IncrRoomAudienceCount(ctx context.Context, roomID string, count int64) error {
	count, err := redis_cli.Client.IncrBy(ctx, roomAudienceCountKey(roomID), count).Result()
	if err != nil {
		return err
	}
	if count == 1 {
		redis_cli.Client.Expire(ctx, roomAudienceCountKey(roomID), expireTime)
	}
	return nil
}

func roomAudienceCountKey(roomID string) string {
	return keyRoomAudienceCountPrefix + roomID
}
