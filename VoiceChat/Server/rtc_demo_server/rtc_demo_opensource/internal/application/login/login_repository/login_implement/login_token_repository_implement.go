package login_implement

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/login/login_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/redis_cli"
)

const (
	fieldUserID      = "user_id"
	fieldCreatedTime = "created_time"
	keyPrefixLogin   = "rtc_demo:login:login_token:"
	TokenExpiration  = 24 * 7 * time.Hour
)

type LoginTokenRepositoryImpl struct{}

func (impl *LoginTokenRepositoryImpl) Save(ctx context.Context, token *login_entity.LoginToken) error {
	logs.CtxInfo(ctx, "redis init login, token: %s, userID: %s", token, token.UserID)
	key := redisKeyLogin(token.Token)

	redis_cli.Client.Expire(ctx, key, TokenExpiration)
	redis_cli.Client.HSet(ctx, key, fieldUserID, token.UserID)
	redis_cli.Client.HSet(ctx, key, fieldCreatedTime, token.CreateTime)
	return nil
}

func (impl *LoginTokenRepositoryImpl) ExistToken(ctx context.Context, token string) (bool, error) {
	key := redisKeyLogin(token)

	val, err := redis_cli.Client.Exists(ctx, key).Result()
	switch err {
	case nil:
		return val == 1, nil
	case redis.Nil:
		return false, nil
	default:
		return false, err
	}
}

func (impl *LoginTokenRepositoryImpl) GetUserID(ctx context.Context, token string) string {
	r := redis_cli.Client.HGet(ctx, redisKeyLogin(token), fieldUserID).Val()
	logs.CtxInfo(ctx, "get login userID: %s", r)

	return r
}

func (impl *LoginTokenRepositoryImpl) GetTokenCreatedAt(ctx context.Context, token string) (int64, error) {
	createdAt, err := redis_cli.Client.HGet(ctx, redisKeyLogin(token), fieldCreatedTime).Int64()
	switch err {
	case nil:
		logs.CtxInfo(ctx, "get token created at: %v", createdAt)
		return createdAt, nil
	case redis.Nil:
		return 0, nil
	default:
		return 0, err
	}
}

func (impl *LoginTokenRepositoryImpl) SetTokenExpiration(ctx context.Context, token string, duration time.Duration) {
	logs.CtxInfo(ctx, "set token expiration: %v", duration)
	redis_cli.Client.Expire(ctx, redisKeyLogin(token), duration)
}

func (impl *LoginTokenRepositoryImpl) GetTokenExpiration(ctx context.Context, token string) time.Duration {
	duration := redis_cli.Client.TTL(ctx, redisKeyLogin(token)).Val()

	logs.CtxInfo(ctx, "get token expiration: %v", duration)
	return duration
}

func redisKeyLogin(token string) string {
	return keyPrefixLogin + token
}
