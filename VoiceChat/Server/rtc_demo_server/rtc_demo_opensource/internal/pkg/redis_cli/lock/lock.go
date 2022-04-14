package lock

import (
	"context"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/redis_cli"
)

const (
	lockTimeout    = 1000 // milliseconds
	tryPeriod      = 10 * time.Millisecond
	lockExpiration = 2 * time.Second
)

const (
	csRoomAssignKey        = "cs:vccontrol:room:assign"
	csRoomLockKeyPrefix    = "cs:vccontrol:room:lock:"
	csLocalUserIDAssignKey = "cs:vccontrol:local:userid:assign"

	ktvRoomSongLockKeyPrefix = "ktv:vscontrol:room_song:lock:"
)

func getDistributedLockKey(roomID string) string {
	return "c:vccontrol:roomlock:" + roomID
}

func LockRoom(ctx context.Context, roomID string) (bool, int64) {
	logs.CtxInfo(ctx, "roomID: %s try to lock", roomID)
	return mustGetLock(ctx, getDistributedLockKey(roomID), 2*time.Second)
}

func UnLockRoom(ctx context.Context, roomID string, lt int64) error {
	logs.CtxInfo(ctx, "roomID: %s try to unlock", roomID)
	return freeLock(ctx, getDistributedLockKey(roomID), lt)
}

func LockCsRoom(ctx context.Context, roomID string) (bool, int64) {
	return mustGetLock(ctx, csRoomLockKeyPrefix+roomID, 2*time.Second)
}

func UnLockCsRoom(ctx context.Context, roomID string, lt int64) error {
	return freeLock(ctx, csRoomLockKeyPrefix+roomID, lt)
}

func LockCsRoomIDAssign(ctx context.Context) (bool, int64) {
	return mustGetLock(ctx, csRoomAssignKey, lockExpiration)
}

func UnLockCsRoomIDAssign(ctx context.Context, lt int64) error {
	return freeLock(ctx, csRoomAssignKey, lt)
}

func LockLocalUserIDAssign(ctx context.Context) (bool, int64) {
	return mustGetLock(ctx, csLocalUserIDAssignKey, lockExpiration)
}

func UnLockLocalUserIDAssign(ctx context.Context, lt int64) error {
	return freeLock(ctx, csLocalUserIDAssignKey, lt)
}

func LockKtvRoomSong(ctx context.Context, roomID string) (bool, int64) {
	return mustGetLock(ctx, ktvRoomSongLockKeyPrefix+roomID, lockExpiration)
}

func UnlockKtvRoomSong(ctx context.Context, roomID string, lt int64) error {
	return freeLock(ctx, ktvRoomSongLockKeyPrefix+roomID, lt)
}

// getLock attemps to get and locking a lock. It returns whether the lock is free and the expiring time of the lock.
func getLock(ctx context.Context, key string) (bool, int64) {
	// lt timestamp microsecond
	lt := time.Now().Add(lockTimeout*time.Millisecond).UnixNano() / 1000
	lock, err := redis_cli.Client.SetNX(ctx, key, lt, lockTimeout*time.Millisecond).Result()
	logs.Error("zhhtest, lock redis_cli error: %v", err)
	if lock {
		return true, lt
	}
	return false, 0
}

// mustGetLock tries to get the lock every 10 ms util obtaining the lock or util timeout.
func mustGetLock(ctx context.Context, key string, timeout time.Duration) (bool, int64) {
	tryTime := timeout / tryPeriod
	for i := 0; i < int(tryTime); i++ {
		l, t := getLock(ctx, key)
		if l {
			return l, t
		}
		time.Sleep(tryPeriod)
	}
	return false, 0
}

// freeLock checks whether the lock is expiring. If not, let it be free.
func freeLock(ctx context.Context, key string, lt int64) error {
	if lt > time.Now().UnixNano()/1000 {
		return redis_cli.Client.Del(ctx, key).Err()
	}
	return nil
}
