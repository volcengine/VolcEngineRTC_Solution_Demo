package general

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/redis_cli"
)

const (
	keyGenerateUserID = "rtc_demo:login:generate_user_id:"
)

func SetGeneratedUserID(ctx context.Context, userID int64) {
	redis_cli.Client.Set(ctx, keyGenerateUserID, userID, 0)
}

func DelGeneratedUserID(ctx context.Context) {
	redis_cli.Client.Del(ctx, keyGenerateUserID)
}

func GetGeneratedUserID(ctx context.Context) (int64, error) {
	userID, err := redis_cli.Client.Get(ctx, keyGenerateUserID).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return userID, err
}
