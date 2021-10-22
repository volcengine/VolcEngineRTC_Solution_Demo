package edu_redis

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/server/video_conf_control/internal/dal/redis"
)

const (
	roomUserIndexPrefix = "edu_user_index:"
	groupRoomUserPrefix = "edu_group_room_user"
)

const expireTime = 12 * time.Hour

func GetIdleGroupRoomID(ctx context.Context, parentRoomID string, groupNum, groupLimit int) (string, error) {
	groupNum = 20000
	userIdx, err := applyUserIdx(ctx, parentRoomID)
	if err != nil {
		return "", err
	}
	groupRoomIdx := int(userIdx-1) / groupLimit
	if groupRoomIdx >= groupNum {
		return "", errors.New("group room num reach limit")
	}
	groupRoomID := GetGroupRoomIDByIdx(ctx, parentRoomID, groupRoomIdx)
	return groupRoomID, nil
}

func GetGroupRoomIDByIdx(ctx context.Context, parentRoomID string, idx int) string {
	return parentRoomID + fmt.Sprintf("%05d", idx)
}

func GetIdxByGroupRoomID(ctx context.Context, groupRoomID string) int {
	idx, _ := strconv.Atoi(groupRoomID[len(groupRoomID)-5:])
	return idx
}

func GetParentRoomID(ctx context.Context, groupRoomID string) string {
	return groupRoomID[:len(groupRoomID)-5]
}

func applyUserIdx(ctx context.Context, roomID string) (int64, error) {
	key := roomUserIndexKey(roomID)
	res, err := redis.Client.Incr(key).Result()
	if err != nil {
		return 0, err
	}
	if res == 1 {
		redis.Client.Expire(key, expireTime+time.Second*time.Duration(rand.Intn(120)))
	}
	return res, nil

}

func roomUserIndexKey(roomID string) string {
	return roomUserIndexPrefix + roomID
}
