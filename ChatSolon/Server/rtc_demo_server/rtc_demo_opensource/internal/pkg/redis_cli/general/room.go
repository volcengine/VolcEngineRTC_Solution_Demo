package general

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/redis_cli"
)

const (
	keyGenerateRoomID = "rtc_demo:generate_room_id:"
)

func SetGeneratedRoomID(ctx context.Context, bizID string, roomID int64) {
	redis_cli.Client.Set(ctx, keyGenerateRoomID+bizID, roomID, 0)
}

func GetGeneratedRoomID(ctx context.Context, bizID string) (int64, error) {
	roomID, err := redis_cli.Client.Get(ctx, keyGenerateRoomID+bizID).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return roomID, err

}
