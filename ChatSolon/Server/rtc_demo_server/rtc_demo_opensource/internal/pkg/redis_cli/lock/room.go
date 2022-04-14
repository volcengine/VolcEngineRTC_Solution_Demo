package lock

import "context"

const (
	generateRoomIDKey = ""
)

func LockGenerateRoomID(ctx context.Context) (bool, int64) {
	return mustGetLock(ctx, generateRoomIDKey, lockExpiration)
}

func UnlockGenerateRoomID(ctx context.Context, lt int64) error {
	return freeLock(ctx, generateRoomIDKey, lt)
}
